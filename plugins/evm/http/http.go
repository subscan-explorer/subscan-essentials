package http

import (
	"github.com/itering/subscan-plugin/router"
	"github.com/itering/subscan/plugins/evm/dao"
	"net/http"
)

func Router() []router.Http {
	srv = &dao.ApiSrv{}
	return []router.Http{
		{"etherscan", etherscanHandle, http.MethodGet},
	}
}
