package storage

import (
	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/itering/subscan/internal/substrate/metadata"
	"github.com/jinzhu/gorm"
)

type Dao interface {
	DB() *gorm.DB
	Redis() *redis.Pool
	RuntimeVersionRaw(int) *metadata.RuntimeRaw
}
