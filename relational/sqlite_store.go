package relational

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"super-descuentos/model"
	"super-descuentos/relational/repository"
)

type SQLStore struct {
	Queries *repository.Queries
	DB      *sql.DB
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		Queries: repository.New(db),
		DB:      db,
	}
}

func (S SQLStore) CreatePost(ctx context.Context, post model.Post) error {

	err := S.Queries.CreatePost(ctx, repository.CreatePostParams{
		ID:           post.ID.String(),
		Title:        post.Title,
		Description:  post.Description,
		Url:          post.Url,
		AuthorID:     post.Author.ID.String(),
		Likes:        int64(post.Likes),
		ExpireTime:   post.ExpireTime,
		CreationTime: post.CreationTime,
	})
	if err != nil {
		_ = fmt.Errorf("hubo un problema al intentar crear el post: %v", err)
		return errors.New("hubo un problema al intentar crear el post")
	}

	return nil
}

func (S SQLStore) DeletePost(ctx context.Context, id uuid.UUID) error {
	result, err := S.Queries.DeletePost(ctx, id.String())
	if err != nil {
		_ = fmt.Errorf("hubo un problema al intentar eliminar el post: %v", err)
		return errors.New("hubo un problema al intentar eliminar el post")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no se encontró un post con el ID: %s", id.String())
	}

	return nil
}

func (S SQLStore) UpdatePost(ctx context.Context, id uuid.UUID, post model.Post) error {
	result, err := S.Queries.UpdatePost(ctx, repository.UpdatePostParams{
		ID:          id.String(),
		Title:       post.Title,
		Description: post.Description,
		Url:         post.Url,
		Likes:       int64(post.Likes),
		ExpireTime:  post.ExpireTime,
	})
	if err != nil {
		_ = fmt.Errorf("hubo un problema al intentar actualizar el post: %v", err)
		return errors.New("hubo un problema al intentar actualizar el post")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no se encontró un post con el ID: %s", id.String())
	}

	return nil
}
func (S SQLStore) GetPost(ctx context.Context, id uuid.UUID) (model.Post, error) {
	post, err := S.Queries.GetPost(ctx, id.String())
	if err != nil {
		_ = fmt.Errorf("hubo un problema al intentar obtener el post: %v", err)
		return model.Post{}, errors.New("hubo un problema al intentar obtener el post")
	}

	user, err := S.Queries.GetUser(ctx, post.AuthorID)
	if err != nil {
		_ = fmt.Errorf("hubo un problema al intentar obtener al autor: %v", err)
		return model.Post{}, errors.New("hubo un problema al intentar obtener al autor")
	}

	return RepositoryPostToModel(post, user), nil
}

func (S SQLStore) GetPosts(ctx context.Context, offset, limit int) ([]model.Post, error) {
	if limit == 0 {
		limit = 10
	}

	posts, err := S.Queries.GetPostsWithAuthor(ctx, repository.GetPostsWithAuthorParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})

	if err != nil {
		_ = fmt.Errorf("hubo un problema al intentar obtener los posts: %v", err)
		return nil, errors.New("hubo un problema al intentar obtener los posts")
	}

	postsWithAuthor := make([]model.Post, 0, limit)
	for _, post := range posts {
		postsWithAuthor = append(postsWithAuthor, RepositoryPostToModel(post.Post, post.User))

	}

	return postsWithAuthor, nil
}
