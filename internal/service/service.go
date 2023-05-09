package service

import (
	"fmt"
	"io/ioutil"
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
	dao dao.IDao
}

// New  a service and return.
func New() (s *Service) {
	websocket.SetEndpoint(util.WSEndPoint)
	d, dbStorage := dao.New()
	s = &Service{dao: d}
	s.initSubRuntimeLatest()
	pluginRegister(dbStorage, d)
	return s
}

func (s *Service) GetDao() dao.IDao {
	return s.dao
}

// Close close the resource.
func (s *Service) Close() {
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
	return ioutil.ReadFile(fmt.Sprintf(util.ConfDir+"/source/%s.json", util.NetworkNode))
}
