package service

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/plugins"
	"github.com/itering/subscan/plugins/storage"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/slog"
)

var (
	subscribeCall      = make(map[string][]plugins.PluginFactory)
	subscribeExtrinsic = make(map[string][]plugins.PluginFactory)
	subscribeEvent     = make(map[string][]plugins.PluginFactory)
)

type extrinsicInfo struct {
	block     *storage.Block
	extrinsic *storage.Extrinsic
	events    []storage.Event
}

type eventInfo struct {
	block     *storage.Block
	extrinsic *storage.Extrinsic
	event     *storage.Event
	fee       decimal.Decimal
}

type callInfo struct {
	block     *storage.Block
	extrinsic *storage.Extrinsic
	call      *storage.Call
	events    []storage.Event
}

type perBlockInfo struct {
	block      *storage.Block
	extrinsics []extrinsicInfo
	events     []eventInfo
	calls      []callInfo
	ready      bool
}

func newPerBlockInfo(block *storage.Block) perBlockInfo {
	return perBlockInfo{
		block: block,
	}
}

type PluginEmitter struct {
	pending       cmap.ConcurrentMap[string, perBlockInfo]
	stop          chan struct{}
	blockComplete chan uint32
	dao           dao.IDao
	s             *Service
}

func NewPluginEmitter(stop chan struct{}, dao dao.IDao, s *Service) PluginEmitter {
	return PluginEmitter{
		pending:       cmap.New[perBlockInfo](),
		stop:          stop,
		blockComplete: make(chan uint32, 1000000),
		dao:           dao,
		s:             s,
	}
}

func (p *PluginEmitter) mutatePerBlock(block *storage.Block, mutate func(*perBlockInfo)) {
	blockKey := fmt.Sprint(block.BlockNum)
	p.pending.Upsert(blockKey, newPerBlockInfo(block), func(exists bool, valInMap perBlockInfo, newVal perBlockInfo) perBlockInfo {
		if !exists {
			mutate(&newVal)
			return newVal
		} else {
			mutate(&valInMap)
			return valInMap
		}
	})
}

func (p *PluginEmitter) noteBlock(block *model.ChainBlock) {
	if err := p.s.NoteChainBlock(block); err != nil {
		slog.Error("NoteChainBlock failed", "error", err)
	}
}

func fallbackProcessStartBlock() int {
	start := os.Getenv("PROCESS_START_BLOCK")
	if start == "" {
		return 0
	}
	parsed, err := strconv.ParseInt(strings.TrimSpace(start), 10, 32)
	if err != nil {
		slog.Error("getProcessStartBlock invalid value", "error", err)
		return 0
	}
	return int(parsed)
}

func (p *PluginEmitter) latestProcessedFromDb() *model.LastProcessedBlock {
	lastProcessed := model.LastProcessedBlock{}
	tx := p.dao.DbBegin()

	tx.Model(&lastProcessed).Select("id", "number").First(&lastProcessed)
	p.dao.DbCommit(tx)
	if lastProcessed.ID == 0 && lastProcessed.Number == 0 {
		return nil
	}
	return &lastProcessed
}

func (p *PluginEmitter) getProcessStartBlock() int {
	num, err := p.dao.GetProcessedBlockNum(context.TODO())
	if err != nil {
		slog.Warn("GetProcessedBlockNum failed", "error", err)
		lastProcessed := p.latestProcessedFromDb()
		slog.Debug("getProcessStartBlock", "lastProcessed", lastProcessed)
		if lastProcessed != nil && lastProcessed.Number > 0 {
			num = lastProcessed.Number
		} else {
			num = fallbackProcessStartBlock()
		}
	}
	return num
}

func (p *PluginEmitter) handleBlock(currentBlock *int64, blockNum uint32) {
	blockNumStr := fmt.Sprint(blockNum)
	if info, ok := p.pending.Get(blockNumStr); ok {
		info.ready = true
		p.pending.Set(blockNumStr, info)
		if next := uint32(*currentBlock + 1); blockNum == next {
			var wait sync.WaitGroup
			nextStr := fmt.Sprint(next)
			for nextInfo, ok := p.pending.Get(nextStr); ok && nextInfo.ready; nextInfo, ok = p.pending.Get(nextStr) {
				slog.Debug("pluginEmitter process block", "blockNum", nextInfo.block.BlockNum)
				wait.Add(1)
				go func() {
					defer wait.Done()
					for _, extrinsic := range nextInfo.extrinsics {
						for _, plugin := range subscribeExtrinsic[extrinsic.extrinsic.CallModule] {
							err := plugin.ProcessExtrinsic(extrinsic.block, extrinsic.extrinsic, extrinsic.events)
							if err != nil {
								slog.Error("ProcessExtrinsic failed", "error", err)
							}
						}
					}
					for _, event := range nextInfo.events {
						for _, plugin := range subscribeEvent[event.event.ModuleId] {
							err := plugin.ProcessEvent(event.block, event.event, event.fee, event.extrinsic)
							if err != nil {
								slog.Error("ProcessEvent failed", "error", err)
							}
						}
					}
					for _, call := range nextInfo.calls {
						for _, plugin := range subscribeCall[call.call.ModuleId] {
							err := plugin.ProcessCall(call.block, call.call, call.events, call.extrinsic)
							if err != nil {
								slog.Error("ProcessCall failed", "error", err)
							}
						}
					}
				}()
				p.pending.Remove(nextStr)
				next++
				nextStr = fmt.Sprint(next)
				wait.Wait()
				*currentBlock += 1
				e := p.dao.SaveProcessedBlockNum(context.TODO(), int(*currentBlock))
				go func() {
					tx := p.dao.DbBegin()
					tx.Model(&model.LastProcessedBlock{}).Where("id = 1").Update("number", currentBlock)
					p.dao.DbCommit(tx)
				}()
				if e != nil {
					slog.Error("SaveProcessedBlockNum failed", "error", e)
				}
			}
		}
	} else {
		slog.Debug("pluginEmitter probably catchUp block (block not in pending)", "blockNum", blockNum)
		block := p.dao.GetBlockByNum(int(blockNum))
		p.noteBlock(block)
		p.handleBlock(currentBlock, blockNum)
	}
}

func (p *PluginEmitter) Run() {
	go func() {
		slog.Debug("pluginEmitter start")
		num := p.getProcessStartBlock()
		currentBlock := int64(num)
		catchUp := p.dao.GetBlocksLaterThan(int(currentBlock))
		slog.Debug("pluginEmitter start", "currentBlock", currentBlock, "catchUp", len(catchUp))
		// populate the `pending` map with the blocks we need to catch up on
		for _, block := range catchUp {
			p.s.blockDone(&block)
		}

		for {
			select {
			case blockNum := <-p.blockComplete:
				slog.Debug("pluginEmitter got block complete", "blockNum", blockNum, "currentBlock", currentBlock)
				p.handleBlock(&currentBlock, blockNum)
			case <-p.stop:
				return
			}
		}
	}()
}

// registered storage
func pluginRegister(ds *dao.DbStorage) {
	for name, plugin := range plugins.RegisteredPlugins {
		db := *ds
		db.Prefix = name
		plugin.InitDao(&db)
		for _, moduleId := range plugin.SubscribeExtrinsic() {
			subscribeExtrinsic[moduleId] = append(subscribeExtrinsic[moduleId], plugin)
		}
		for _, moduleId := range plugin.SubscribeEvent() {
			subscribeEvent[moduleId] = append(subscribeEvent[moduleId], plugin)
		}
		for _, moduleId := range plugin.SubscribeCall() {
			subscribeCall[moduleId] = append(subscribeCall[moduleId], plugin)
		}
	}
}

func (s *Service) blockDone(block *model.ChainBlock) {
	slog.Debug("blockDone", "block", block.BlockNum)
	s.pluginEmitter.blockComplete <- uint32(block.BlockNum)
}

// after event created, emit event data to subscribe plugins
func (s *Service) emitEvent(block *model.ChainBlock, event *model.ChainEvent, fee decimal.Decimal, extrinsic *model.ChainExtrinsic) {
	pBlock := block.AsPlugin()
	s.pluginEmitter.mutatePerBlock(pBlock, func(info *perBlockInfo) {
		var ext *storage.Extrinsic
		if extrinsic != nil {
			ext = extrinsic.AsPlugin()
		}
		info.events = append(info.events, eventInfo{
			block:     pBlock,
			extrinsic: ext,
			event:     event.AsPlugin(),
			fee:       fee,
		})
	})
}

// after extrinsic created, emit extrinsic data to subscribe plugins
func (s *Service) emitExtrinsic(block *model.ChainBlock, extrinsic *model.ChainExtrinsic, events []model.ChainEvent) {
	block.BlockTimestamp = extrinsic.BlockTimestamp

	pBlock := block.AsPlugin()
	s.pluginEmitter.mutatePerBlock(pBlock, func(info *perBlockInfo) {
		info.extrinsics = append(info.extrinsics, extrinsicInfo{
			block:     pBlock,
			extrinsic: extrinsic.AsPlugin(),
			events:    model.MapAsPlugin[*storage.Event](events),
		})
	})
}

func (s *Service) emitCall(block *model.ChainBlock, call *model.ChainCall, events []model.ChainEvent, extrinsic *model.ChainExtrinsic) {
	slog.Debug("emit call", "subscribeCall", subscribeCall)
	pBlock := block.AsPlugin()
	s.pluginEmitter.mutatePerBlock(pBlock, func(info *perBlockInfo) {
		info.calls = append(info.calls, callInfo{
			block:     pBlock,
			extrinsic: extrinsic.AsPlugin(),
			call:      call.AsPlugin(),
			events:    model.MapAsPlugin[*storage.Event](events),
		})
	})
}
