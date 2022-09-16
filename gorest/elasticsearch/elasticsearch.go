package elasticsearch

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticSearch struct {
	client *elasticsearch.Client
	index  string
	alias  string
}

func New(addresses []string, username string, password string) (*ElasticSearch, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
		Username:  username,
		Password:  password,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	res, err := es.Info()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("Error: %s", res.String())
	}

	return &ElasticSearch{
		client: es,
	}, nil
}

func (e *ElasticSearch) CreateIndex(index string) error {
	e.index = index
	e.alias = index + "_alias"

	res, err := e.client.Indices.Exists([]string{e.index})
	if err != nil {
		return fmt.Errorf("Cannot check index existence: %w", err)
	}
	if res.StatusCode == 200 {
		return nil
	}
	if res.StatusCode != 404 {
		return fmt.Errorf("Error in index existence response: %s", res.String())
	}

	res, err = e.client.Indices.Create(e.index)
	if err != nil {
		return fmt.Errorf("Cannot create index: %w", err)
	}
	if res.IsError() {
		return fmt.Errorf("Error in index creation response: %s", res.String())
	}

	res, err = e.client.Indices.PutAlias([]string{e.index}, e.alias)
	if err != nil {
		return fmt.Errorf("Cannot create index alias: %w", err)
	}
	if res.IsError() {
		return fmt.Errorf("Error in index alias creation response: %s", res.String())
	}

	return nil
}

type document struct {
	Source interface{} `json:"_source"`
}
