package relational

import (
	"super-descuentos/model"
	"super-descuentos/relational/repository"

	"github.com/google/uuid"
)

func RepositoryPostToModel(post repository.Post, user repository.User) model.Post {
	return model.Post{
		ID:           uuid.MustParse(post.ID),
		Title:        post.Title,
		Description:  post.Description,
		Url:          post.Url,
		Author:       model.User{ID: uuid.MustParse(user.ID), Name: user.Name, Email: user.Email},
		Likes:        int(post.Likes),
		ExpireTime:   post.ExpireTime,
		CreationTime: post.CreationTime,
	}
}

func RepositoryAuthorToModel(author repository.User) model.User {
	return model.User{
		ID:    uuid.MustParse(author.ID),
		Name:  author.Name,
		Email: author.Email,
	}
}
