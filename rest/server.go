package rest

import (
	"net/http"
	"super-descuentos/store"
)

type Server struct {
	store store.Store
	http.Handler
}

func NewServer(store store.Store) *Server {
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
