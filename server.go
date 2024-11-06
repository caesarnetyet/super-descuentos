package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
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

func (s *Server) validateUUID(idStr string) (uuid.UUID, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, ErrInvalidID
	}
	return id, nil
}

func (s *Server) jsonWithErrors(w http.ResponseWriter, v interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func (s *Server) sendErrorMessage(w http.ResponseWriter, err error, code int) {
	s.jsonWithErrors(w, map[string]string{"message": err.Error()}, code)
}
