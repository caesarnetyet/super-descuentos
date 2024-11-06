package main

import (
	"github.com/google/uuid"
	"net/http"
)

type Store interface {
	CreatePost(post Post) error
	DeletePost(id uuid.UUID) error
	UpdatePost(id uuid.UUID, post Post) error
	GetPost(id uuid.UUID) (Post, error)
	GetPosts() ([]Post, error)
}

type Server struct {
	store Store
	http.Handler
}

func NewServer(store Store) *Server {
	s := new(Server)
	s.store = store
	router := http.NewServeMux()
	router.HandleFunc("GET /posts", s.handlePosts)
	router.HandleFunc("GET /posts/{id}", s.handlePost)
	router.HandleFunc("POST /posts", s.handleCreatePost)
	router.HandleFunc("DELETE /posts/{id}", s.handleDeletePost)
	router.HandleFunc("PUT /posts/{id}", s.handleUpdatePost)
	s.Handler = router
	return s
}
