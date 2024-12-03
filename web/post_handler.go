package web

import (
	"net/http"
	"super-descuentos/components"
	"super-descuentos/model"
	"super-descuentos/utils"

	"github.com/google/uuid"
)

func (s *Server) handlePostForm(w http.ResponseWriter, r *http.Request) {
	authors, err := s.store.GetAuthors(r.Context(), 0, 100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = components.Layout("Posts", components.PostsPage(authors)).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type HandleCreatePostFormRequest struct {
	Title       string `json:"title" form:"title"`
	Content     string `json:"content" form:"content"`
	AuthorEmail string `json:"author_email" form:"author_email"`
	Url         string `json:"url" form:"url"`
}

func (r HandleCreatePostFormRequest) Validate() model.ValidationErrors {
	var errs model.ValidationErrors

	if r.Title == "" {
		errs = append(errs, model.ValidationError{
			Field:   "title",
			Message: "el título es requerido",
		})
	}

	if r.Content == "" {
		errs = append(errs, model.ValidationError{
			Field:   "content",
			Message: "el contenido es requerido",
		})
	}

	if r.AuthorEmail == "" {
		errs = append(errs, model.ValidationError{
			Field:   "author_email",
			Message: "el correo electrónico del autor es requerido",
		})
	}

	if r.Url == "" {
		errs = append(errs, model.ValidationError{
			Field:   "url",
			Message: "la URL es requerida",
		})
	}

	return errs
}

func (s *Server) handleCreatePostForm(w http.ResponseWriter, r *http.Request) {

	request, err := model.DecodeAndValidate[HandleCreatePostFormRequest](r)
	if err != nil {
		utils.HandleErrorResponse(w, err)
		return
	}

	author, err := s.store.GetAuthorByEmail(r.Context(), request.AuthorEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post := model.Post{
		ID:          uuid.New(),
		Title:       request.Title,
		Description: request.Content,
		Url:         request.Url,
		Author:      author,
	}

	errs := post.Validate()
	if len(errs) > 0 {
		utils.HandleErrorResponse(w, errs)
		return
	}

	err = s.store.CreatePost(r.Context(), post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
