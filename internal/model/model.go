package model

import (
	"subscan-end/utiles/es"
)

var esClient *es.EsClient

func InitEsClient() {
	if esClient == nil {
		esClient, _ = es.NewEsClient()
	}
}

