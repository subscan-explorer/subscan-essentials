package es

import (
	"context"
	"github.com/olivere/elastic"
)

func (c *EsClient) Insert(index, id string, bodyJson interface{}) error {
	_, err := c.Client.Index().Index(index).Id(id).BodyJson(bodyJson).Do(context.Background())
	return err
}

func (c *EsClient) Update(index, id string, bodyJson interface{}) error {
	_, err := c.Client.Update().Index(index).Id(id).Doc(bodyJson).Do(context.Background())
	return err
}

func (c *EsClient) Delete(index, id string) error {
	_, err := c.Client.Delete().Index(index).Id(id).Do(context.Background())
	return err
}

func (c *EsClient) GetById(index, id string) (interface{}, error) {
	got, err := c.Client.Get().Index(index).Id(id).Do(context.Background())
	if err == nil {
		_, _ = c.Client.Refresh().Index(index).Do(context.Background())
	} else {
		return nil, err
	}
	return got.Source, err
}

func (c *EsClient) GetWhere(index string, termQuery elastic.Query, sortField string, asc bool, offset, size int) ([]interface{}, error) {
	searchResult, err := c.Client.Search().Index(index).Query(termQuery).Sort(sortField, asc).From(offset).Size(size).Pretty(true).Do(context.Background())
	if searchResult == nil || err != nil {
		return nil, err
	}
	var result []interface{}
	if searchResult.TotalHits() > 0 {
		for _, hit := range searchResult.Hits.Hits {
			result = append(result, hit.Source)
		}
	}
	return result, err
}
