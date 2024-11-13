package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"super-descuentos/errs"
	"super-descuentos/model"
	"super-descuentos/validator"
	"time"

	"github.com/google/uuid"
)

func (s *Server) handlePosts(w http.ResponseWriter, r *http.Request) {
	posts, err := s.store.GetPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (s *Server) handlePost(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := s.validateUUID(idStr)
	if err != nil {
		s.jsonWithErrors(w, map[string]string{"message": err.Error()}, http.StatusBadRequest)
		return
	}

	post, err := s.store.GetPost(id)
	if err != nil {
		if errors.Is(err, errs.ErrPostNotFound) {
			s.sendErrorMessage(w, err, http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (s *Server) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	post, err := validator.DecodeAndValidate[model.Post](r)
	if err != nil {
		switch e := err.(type) {
		case validator.ValidationError:
			s.sendErrorMessage(w, errors.New(e.Message), http.StatusBadRequest)
		case validator.ValidationErrors:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": e,
			})
		default:
			http.Error(w, "error interno del servidor", http.StatusInternalServerError)
		}
		return
	}

	post.ID = uuid.New()
	post.CreationTime = time.Now()

	if err := s.store.CreatePost(post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (s *Server) handleDeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := s.validateUUID(idStr)
	if err != nil {
		s.jsonWithErrors(w, map[string]string{"message": err.Error()}, http.StatusBadRequest)
		return
	}

	if err := s.store.DeletePost(id); err != nil {
		if errors.Is(err, errs.ErrPostNotFound) {
			s.sendErrorMessage(w, err, http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleUpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := s.validateUUID(idStr)
	if err != nil {
		s.jsonWithErrors(w, map[string]string{"message": err.Error()}, http.StatusBadRequest)
		return
	}

	post, err := validator.DecodeAndValidate[model.Post](r)
	if err != nil {
		switch e := err.(type) {
		case validator.ValidationError:
			s.sendErrorMessage(w, errors.New(e.Message), http.StatusBadRequest)
		case validator.ValidationErrors:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": e,
			})
		default:
			http.Error(w, "error interno del servidor", http.StatusInternalServerError)

		}
	}

	if err := s.store.UpdatePost(id, post); err != nil {
		if errors.Is(err, errs.ErrPostNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
