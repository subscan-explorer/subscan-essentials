package service

import (
	"fmt"
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
}

func NewPluginEmitter(stop chan struct{}) PluginEmitter {
	return PluginEmitter{
		pending:       cmap.New[perBlockInfo](),
		stop:          stop,
		blockComplete: make(chan uint32),
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

func (p *PluginEmitter) Run() {
	go func() {
		currentBlock := int64(0)
		slog.Debug("pluginEmitter start")
		for {
			select {
			case blockNum := <-p.blockComplete:
				slog.Debug("pluginEmitter got block complete", "blockNum", blockNum, "currentBlock", currentBlock)
				blockNumStr := fmt.Sprint(blockNum)
				if info, ok := p.pending.Get(blockNumStr); ok {
					info.ready = true
					p.pending.Set(blockNumStr, info)
					if next := uint32(currentBlock + 1); blockNum == next {
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
						}
						currentBlock = int64(next) - 1
					}
				}
			case <-p.stop:
				return
			}
		}
	}()
}

// registered storage
func pluginRegister(ds *dao.DbStorage, dd *dao.Dao) {
	for name, plugin := range plugins.RegisteredPlugins {
		db := *ds
		db.Prefix = name
		plugin.InitDao(&db, dd)
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
		info.events = append(info.events, eventInfo{
			block:     pBlock,
			extrinsic: extrinsic.AsPlugin(),
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
			block:     block.AsPlugin(),
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
