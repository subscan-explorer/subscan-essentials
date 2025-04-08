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
)

// Service
type Service struct {
	dao       dao.IDao
	dbStorage *dao.DbStorage
}

// New  a service and return.
func New() (s *Service) {
	websocket.SetEndpoint(util.WSEndPoint)
	d, dbStorage, pool := dao.New()
	s = &Service{dao: d, dbStorage: dbStorage}
	s.initSubRuntimeLatest()
	pluginRegister(dbStorage, pool)
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
		if data, err := readTypeRegistry(); err == nil {
			substrate.RegCustomTypes(data)
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
