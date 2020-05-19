package es

import (
	"errors"
	"github.com/olivere/elastic"
	"log"
	"os"
	"subscan-end/utiles"
	"time"
)

type EsClient struct {
	Client     *elastic.Client
	InfoLogger *log.Logger
	ErrLogger  *log.Logger
}

func NewEsClient() (*EsClient, error) {
	if utiles.GetEnv("ES_ENABLE", "true") == "false" {
		return nil, errors.New("es not enable")
	}
	InfoLogger := log.New(os.Stderr, "ELASTIC_INFO ", log.LstdFlags)
	ErrLogger := log.New(os.Stdout, "ELASTIC_ERR ", log.LstdFlags)
	client, err := elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200"),
		elastic.SetSniff(false), // Cluster Nodes
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetGzip(true),
		elastic.SetErrorLog(ErrLogger),
		elastic.SetInfoLog(InfoLogger),
	)
	if err != nil {
		return nil, err
	}
	esClient := EsClient{Client: client, InfoLogger: InfoLogger, ErrLogger: InfoLogger}
	return &esClient, nil
}
