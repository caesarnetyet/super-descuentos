package relational

import (
	"github.com/google/uuid"
	"super-descuentos/model"
	"super-descuentos/relational/repository"
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