package service

import (
	"fmt"
	"os"
	"strings"

	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/websocket"
	"golang.org/x/exp/slog"
)

// Service
type Service struct {
	dao           dao.IDao
	pluginEmitter PluginEmitter
	ReadOnlyService
}

// New  a service and return.
func New(stop chan struct{}) (s *Service) {
	d, dbStorage := dao.New(true)
	pluginRegister(dbStorage)
	return newWithDao(stop, d)
}

func newWithDao(stop chan struct{}, dao dao.IDao) *Service {
	websocket.SetEndpoint(util.WSEndPoint)
	s := &Service{dao: dao, ReadOnlyService: *readOnlyWithDao(dao)}
	s.initSubRuntimeLatest()
	s.pluginEmitter = NewPluginEmitter(stop)
	return s
}

func (s *Service) Run() {
	s.pluginEmitter.Run()
}

func (s *Service) GetDao() dao.IDao {
	return s.dao
}

// Close close the resource.
func (s *ReadOnlyService) Close() {
	s.dao.Close()
}

func (s *Service) initSubRuntimeLatest() {
	// reg network custom type
	defer func() {
		go s.unknownToken()
		if c, err := readTypeRegistry(); err == nil {
			substrate.RegCustomTypes(c)
			if unknown := metadata.Decoder.CheckRegistry(); len(unknown) > 0 {
				slog.Warn("Found unknown type %s", strings.Join(unknown, ", "))
			}
		} else {
			if os.Getenv("TEST_MOD") != "true" {
				panic(err)
			}
		}
	}()

	// find db
	if recent := s.dao.RuntimeVersionRecent(); recent != nil && strings.HasPrefix(recent.RawData, "0x") {
		metadata.Latest(&metadata.RuntimeRaw{Spec: recent.SpecVersion, Raw: recent.RawData})
		return
	}
	// find metadata for blockChain
	if raw := s.regCodecMetadata(); strings.HasPrefix(raw, "0x") {
		metadata.Latest(&metadata.RuntimeRaw{Spec: 1, Raw: raw})
		return
	}
	panic("Can not find chain metadata, please check network")
}

// read custom registry from local or remote
func readTypeRegistry() ([]byte, error) {
	return os.ReadFile(fmt.Sprintf(util.ConfDir+"/source/%s.json", util.NetworkNode))
}
