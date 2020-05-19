package es

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
)

type Index struct {
	Settings struct {
		NumberOfShards   int `json:"number_of_shards"`
		NumberOfReplicas int `json:"number_of_replicas"`
	} `json:"settings"`
	Mappings Mappings `json:"mappings"`
}

type Mappings struct {
	Doc Doc `json:"doc"`
}

type Doc struct {
	Properties map[string]map[string]string `json:"properties"`
}

func NewIndexTemplate() *Index {
	i := Index{}
	i.Settings.NumberOfShards = 1
	i.Settings.NumberOfReplicas = 0
	i.Mappings.Doc.Properties = make(map[string]map[string]string)
	return &i
}

func (index *Index) InjectIndex(s interface{}) {
	t := reflect.TypeOf(s)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get("es") != "" && field.Tag.Get("json") != "" {
			indexTag := strings.Split(field.Tag.Get("es"), ",")
			if len(indexTag) > 0 {
				thisIndex := make(map[string]string)
				for _, tag := range indexTag {
					tagArr := strings.Split(tag, ":")
					if len(tagArr) == 2 {
						thisIndex[tagArr[0]] = tagArr[1]
					}
				}
				index.Mappings.Doc.Properties[field.Tag.Get("json")] = thisIndex
			}
		}
	}
}

func (c *EsClient) CreateIndex(i *Index, indexName string) (err error) {
	bIndex, _ := json.Marshal(i)
	if exist, _ := c.Client.IndexExists(indexName).Do(context.Background()); !exist {
		_, err = c.Client.CreateIndex(indexName).Body(string(bIndex)).Do(context.Background())
	}
	return err
}

func (c *EsClient) DelIndex(indexName string) error {
	var (
		err   error
		exist bool
	)
	if exist, err = c.Client.IndexExists(indexName).Do(context.Background()); exist {
		_, err = c.Client.DeleteIndex(indexName).Do(context.Background())
	}
	return err
}
