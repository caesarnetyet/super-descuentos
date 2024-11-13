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

	// Intentar cargar templates, pero no fallar si no existen
	if templates, err := template.ParseGlob("templates/**/*.html"); err == nil {
		s.templates = templates
	}

	router := http.NewServeMux()

	// Rutas de la API
	router.HandleFunc("GET /api/posts", s.handlePosts)
	router.HandleFunc("GET /api/posts/{id}", s.handlePost)
	router.HandleFunc("POST /api/posts", s.handleCreatePost)
	router.HandleFunc("DELETE /api/posts/{id}", s.handleDeletePost)
	router.HandleFunc("PUT /api/posts/{id}", s.handleUpdatePost)

	// Rutas para las vistas web (solo si hay templates)
	if s.templates != nil {
		router.HandleFunc("GET /web/posts", s.handlePostsPage)
		router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				http.Redirect(w, r, "/web/posts", http.StatusFound)
				return
			}
		})

		// Archivos estáticos
		router.HandleFunc("GET /static/{file}", func(w http.ResponseWriter, r *http.Request) {
			file := r.PathValue("file")
			if file == "script.js" || file == "style.css" {
				http.ServeFile(w, r, filepath.Join("templates/posts", file))
			}
		})
	}

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
