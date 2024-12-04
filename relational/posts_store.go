package relational

import (
	"context"
	"errors"
	"fmt"
	"super-descuentos/errs"
	"super-descuentos/model"
	"super-descuentos/relational/repository"

	"github.com/google/uuid"
)

func (store *SQLStore) CreatePost(ctx context.Context, post model.Post) error {
	_, err := store.Queries.GetUser(ctx, post.Author.ID.String())

	if err != nil {
		_ = fmt.Errorf("hubo un problema al intentar obtener al autor: %v", err)
		return errs.ErrAuthorNotFound
	}
	err = store.Queries.CreatePost(ctx, repository.CreatePostParams{
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

func (store *SQLStore) DeletePost(ctx context.Context, id uuid.UUID) error {
	result, err := store.Queries.DeletePost(ctx, id.String())
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

func (store *SQLStore) UpdatePost(ctx context.Context, id uuid.UUID, post model.Post) error {
	result, err := store.Queries.UpdatePost(ctx, repository.UpdatePostParams{
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
func (store *SQLStore) GetPost(ctx context.Context, id uuid.UUID) (model.Post, error) {
	post, err := store.Queries.GetPost(ctx, id.String())
	if err != nil {
		_ = fmt.Errorf("hubo un problema al intentar obtener el post: %v", err)
		return model.Post{}, errors.New("hubo un problema al intentar obtener el post")
	}

	user, err := store.Queries.GetUser(ctx, post.AuthorID)
	if err != nil {
		_ = fmt.Errorf("hubo un problema al intentar obtener al autor: %v", err)
		return model.Post{}, errors.New("hubo un problema al intentar obtener al autor")
	}

	return RepositoryPostToModel(post, user), nil
}

func (store *SQLStore) GetPosts(ctx context.Context, offset, limit int) ([]model.Post, error) {
	if limit == 0 {
		limit = 10
	}

	posts, err := store.Queries.GetPostsWithAuthor(ctx, repository.GetPostsWithAuthorParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})

	if err != nil {
		fmt.Printf("hubo un problema al intentar obtener los posts: %v", err)
		return nil, errors.New("hubo un problema al intentar obtener los posts")
	}

	postsWithAuthor := make([]model.Post, 0, limit)
	for _, post := range posts {
		postsWithAuthor = append(postsWithAuthor, RepositoryPostToModel(post.Post, post.User))

	}

	return postsWithAuthor, nil
}

func (store *SQLStore) CreateAuthor(ctx context.Context, author model.User) error {
	err := store.Queries.CreateUser(ctx, repository.CreateUserParams{
		ID:    author.ID.String(),
		Name:  author.Name,
		Email: author.Email,
	})
	if err != nil {
		fmt.Printf("hubo un problema al intentar crear el autor: %v \n", err)
		return errors.New("hubo un problema al intentar crear el autor")
	}

	return nil
}
