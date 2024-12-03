package rest

import (
	"net/http"
	"super-descuentos/store"
)

type Server struct {
	store store.Store
	http.Handler
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func NewServer(store store.Store) *Server {
	s := new(Server)
	s.store = store

	router := http.NewServeMux()

	// Rutas de la API
	router.HandleFunc("GET /posts", s.handlePosts)
	router.HandleFunc("GET /posts/{id}", s.handlePost)
	router.HandleFunc("POST /posts", s.handleCreatePost)
	router.HandleFunc("DELETE /posts/{id}", s.handleDeletePost)
	router.HandleFunc("PUT /posts/{id}", s.handleUpdatePost)

	s.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCORS(router).ServeHTTP(w, r)
	})

	return s
}
