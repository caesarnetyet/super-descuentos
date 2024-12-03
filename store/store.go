package store

import (
	"context"
	"super-descuentos/model"

	"github.com/google/uuid"
)

type Store interface {
	CreatePost(ctx context.Context, post model.Post) error
	DeletePost(ctx context.Context, id uuid.UUID) error
	UpdatePost(ctx context.Context, id uuid.UUID, post model.Post) error
	GetPost(ctx context.Context, id uuid.UUID) (model.Post, error)
	GetPosts(ctx context.Context, offset, limit int) ([]model.Post, error)
	GetAuthorByEmail(ctx context.Context, email string) (model.User, error)
	GetAuthors(ctx context.Context, offset, limit int) ([]model.User, error)
	CreateAuthor(ctx context.Context, author model.User) error
}
