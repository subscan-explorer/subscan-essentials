package service

import (
	"fmt"
	"github.com/itering/scale.go/source"
	"github.com/itering/scale.go/types"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/internal/service/scan"
	"github.com/itering/subscan/lib/substrate/metadata"
	"github.com/itering/subscan/plugins"
	"github.com/itering/subscan/util"
	"io/ioutil"
	"strings"
)

// Service service.
type Service struct {
	dao *dao.Dao
}

// New new a service and return.
func New() (s *Service) {
	s = &Service{
		dao: dao.New(),
	}

	s.Migration()
	s.initSubRuntimeLatest()
	// load plugins
	for _, plugin := range plugins.RegisteredPlugins {
		p := plugin()
		p.InitDao(s.dao)
	}
	return s
}

func (s *Service) NewScan() *scan.Service {
	return scan.New(s.dao)
}

// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
}

func (s *Service) Migration() {
	s.dao.Migration()
}

func (s *Service) initSubRuntimeLatest() {
	// reg network custom type
	defer func() {
		go s.UnknownToken()
		c, err := ioutil.ReadFile(fmt.Sprintf("../configs/source/%s.json", util.NetworkNode))
		if err == nil {
			types.RegCustomTypes(source.LoadTypeRegistry(c))
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
	panic("can not find chain metadata")
}
