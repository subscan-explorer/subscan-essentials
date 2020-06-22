package service

import (
	"fmt"
	"github.com/itering/scale.go/source"
	"github.com/itering/scale.go/types"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/internal/service/scan"
	"github.com/itering/subscan/internal/substrate/metadata"
	"github.com/itering/subscan/internal/util"
	"io/ioutil"
	"strings"
)

// Service service.
type Service struct {
	Dao *dao.Dao
}

// New new a service and return.
func New() (s *Service) {
	s = &Service{
		Dao: dao.New(),
	}

	s.Migration()
	s.initSubRuntimeLatest()
	return s
}

func (s *Service) NewScan() *scan.Service {
	return scan.New(s.Dao)
}

// Close close the resource.
func (s *Service) Close() {
	s.Dao.Close()
}

func (s *Service) Migration() {
	s.Dao.Migration()
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
	if recent := s.Dao.RuntimeVersionRecent(); recent != nil && strings.HasPrefix(recent.RawData, "0x") {
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
