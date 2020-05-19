package es

import (
	"errors"
	"github.com/itering/subscan/util"
	"github.com/olivere/elastic"
	"log"
	"os"
	"time"
)

type EsClient struct {
	Client     *elastic.Client
	InfoLogger *log.Logger
	ErrLogger  *log.Logger
}

func NewEsClient() (*EsClient, error) {
	if util.GetEnv("ES_ENABLE", "true") == "false" {
		return nil, errors.New("es not enable")
	}
	InfoLogger := log.New(os.Stderr, "ELASTIC_INFO ", log.LstdFlags)
	ErrLogger := log.New(os.Stdout, "ELASTIC_ERR ", log.LstdFlags)
	client, err := elastic.NewClient(
		// TODO
		//
		// Default URL
		elastic.SetURL(""),
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
