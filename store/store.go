package store

import (
	"context"
	"github.com/google/uuid"
	"super-descuentos/model"
)

type Store interface {
	CreatePost(ctx context.Context, post model.Post) error
	DeletePost(ctx context.Context, id uuid.UUID) error
	UpdatePost(ctx context.Context, id uuid.UUID, post model.Post) error
	GetPost(ctx context.Context, id uuid.UUID) (model.Post, error)
	GetPosts(ctx context.Context, offset, limit int) ([]model.Post, error)
}
