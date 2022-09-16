package post

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type service struct {
	storage post_storage.PostStorer
}

func (s service) create(ctx context.Context, req createRequest) (createResponse, error) {
	id := uuid.New().String()
	cr := time.Now().UTC()

	doc := post_storage.Post{
		ID:        id,
		Title:     req.Title,
		Text:      req.Text,
		Tags:      req.Tags,
		CreatedAt: &cr,
	}

	if err := s.storage.Insert(ctx, doc); err != nil {
		return createResponse{}, err
	}
	return createResponse{ID: id}, nil
}
