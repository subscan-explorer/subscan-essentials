package scan

import (
	"context"
	"fmt"
	"github.com/itering/subscan/internal/util"
)

func (s *Service) CacheFunc(call func() interface{}, refresh bool) []byte {
	ctx := context.TODO()

	caller := util.CallerName()
	key := fmt.Sprintf("%s:scan:%s", util.NetworkNode, caller)

	if !refresh {
		if cacheData := s.dao.GetCacheBytes(ctx, key); cacheData != nil {
			return cacheData
		}
	}

	value := call()
	if value == nil {
		return nil
	}

	_ = s.dao.SetCache(context.TODO(), key, value, 3600)
	b, _ := json.Marshal(value)
	return b
}
