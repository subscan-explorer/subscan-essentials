// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package di

import (
	// "github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/internal/server/http"
	"github.com/itering/subscan/internal/service"

	"github.com/google/wire"
)

//go:generate kratos t wire
func InitApp() (*App, func(), error) {
	panic(wire.Build(service.Provider, http.New, NewApp))
}
