package rest

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"super-descuentos/errs"
	"super-descuentos/model"

	"github.com/google/uuid"
)

type Store interface {
	CreatePost(post model.Post) error
	DeletePost(id uuid.UUID) error
	UpdatePost(id uuid.UUID, post model.Post) error
	GetPost(id uuid.UUID) (model.Post, error)
	GetPosts() ([]model.Post, error)
}

type Server struct {
	store     Store
	templates *template.Template
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

func NewServer(store Store) *Server {
	s := new(Server)
	s.store = store

	templates, err := template.ParseGlob(filepath.Join("templates", "*.html"))
	if err != nil {
		panic(err)
	}
	s.templates = templates

	router := http.NewServeMux()

	// Rutas de la API
	router.HandleFunc("GET /api/posts", s.handlePosts)
	router.HandleFunc("GET /api/posts/{id}", s.handlePost)
	router.HandleFunc("POST /api/posts", s.handleCreatePost)
	router.HandleFunc("DELETE /api/posts/{id}", s.handleDeletePost)
	router.HandleFunc("PUT /api/posts/{id}", s.handleUpdatePost)

	// Rutas para las vistas web
	router.HandleFunc("GET /web/posts", s.handlePostsPage)
	router.HandleFunc("GET /", s.handleHome)

	// Aplicar CORS a todas las rutas /api/
	s.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			enableCORS(router).ServeHTTP(w, r)
		} else {
			router.ServeHTTP(w, r)
		}
	})

	return s
}

func (s *Server) validateUUID(idStr string) (uuid.UUID, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, errs.ErrInvalidID
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
