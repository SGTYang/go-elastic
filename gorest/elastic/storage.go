package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var _ PostStorer = PostStorage{}

type PostStorer interface {
	Insert(ctx context.Context, post Post) error
}

type Post struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Text      string     `json:"text"`
	Tags      []string   `json:"tags"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

type PostStorage struct {
	elastic ElasticSearch
	timeout time.Duration
}

func NewPostStorage(elastic ElasticSearch) (PostStorage, error) {
	return PostStorage{
		elastic: elastic,
		timeout: time.Second * 10,
	}, nil
}

func (p PostStorage) Insert(ctx context.Context, post Post) error {
	bdy, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("Insert: marshall: %w", err)
	}
	req := esapi.CreateRequest{
		Index:      p.elastic.alias,
		DocumentID: post.ID,
		Body:       bytes.NewReader(bdy),
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()
	res, err := req.Do(ctx, p.elastic.client)
	if err != nil {
		return fmt.Errorf("Insert: request: %w", err)
	}
	if res.IsError() {
		return fmt.Errorf("Insert: response: %s", res.String())
	}
	return nil
}
