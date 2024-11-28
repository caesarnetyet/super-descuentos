package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"super-descuentos/errs"
	"super-descuentos/model"
	"super-descuentos/utils"
	"time"

	"github.com/google/uuid"
)

func (s *Server) handlePosts(w http.ResponseWriter, r *http.Request) {
	offset, limit, err := utils.GetOffsetAndLimit(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	posts, err := s.store.GetPosts(r.Context(), offset, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (s *Server) handlePost(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := utils.ValidateUUID(idStr)
	if err != nil {
		utils.JsonWithErrors(w, map[string]string{"message": err.Error()}, http.StatusBadRequest)
		return
	}

	post, err := s.store.GetPost(r.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrPostNotFound) {
			utils.SendErrorMessage(w, err, http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (s *Server) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	post, err := model.DecodeAndValidate[model.Post](r)
	if err != nil {
		switch e := err.(type) {
		case model.ValidationError:
			utils.SendErrorMessage(w, errors.New(e.Message), http.StatusBadRequest)
		case model.ValidationErrors:
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

	if err := s.store.CreatePost(r.Context(), post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (s *Server) handleDeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := utils.ValidateUUID(idStr)
	if err != nil {
		utils.JsonWithErrors(w, map[string]string{"message": err.Error()}, http.StatusBadRequest)
		return
	}

	if err := s.store.DeletePost(r.Context(), id); err != nil {
		if errors.Is(err, errs.ErrPostNotFound) {
			utils.SendErrorMessage(w, err, http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleUpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := utils.ValidateUUID(idStr)
	if err != nil {
		utils.JsonWithErrors(w, map[string]string{"message": err.Error()}, http.StatusBadRequest)
		return
	}

	post, err := model.DecodeAndValidate[model.Post](r)
	if err != nil {
		switch e := err.(type) {
		case model.ValidationError:
			utils.SendErrorMessage(w, errors.New(e.Message), http.StatusBadRequest)
		case model.ValidationErrors:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": e,
			})
		default:
			http.Error(w, "error interno del servidor", http.StatusInternalServerError)

		}
	}

	if err := s.store.UpdatePost(r.Context(), id, post); err != nil {
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
