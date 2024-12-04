package web

import (
	"github.com/google/uuid"
	"net/http"
	"super-descuentos/components"
	"super-descuentos/errs"
	"super-descuentos/model"
	"super-descuentos/utils"
)

func (s *Server) handleAuthorForm(w http.ResponseWriter, r *http.Request) {
	offset, limit, err := utils.GetOffsetAndLimit(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authors, err := s.store.GetAuthors(r.Context(), offset, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	component := components.Layout("Authors", components.AuthorsPage(authors))
	err = component.Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

type HandleCreateAuthorFormRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (r HandleCreateAuthorFormRequest) Validate() model.ValidationErrors {
	var problems model.ValidationErrors

	if r.Name == "" {
		problems = append(problems, model.ValidationError{
			Field:   "name",
			Message: "el nombre es requerido",
		})
	}

	if r.Email == "" {
		problems = append(problems, model.ValidationError{
			Field:   "email",
			Message: "el correo electrónico es requerido",
		})
	}

	return problems
}

func (s *Server) handleCreateAuthorForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, errs.ErrInvalidFormData.Error(), http.StatusBadRequest)
		return
	}

	author := model.User{
		ID:    uuid.New(),
		Name:  r.FormValue("name"),
		Email: r.FormValue("email"),
	}

	if problems := author.Validate(); len(problems) > 0 {
		http.Error(w, problems.Error(), http.StatusBadRequest)
		return
	}

	err = s.store.CreateAuthor(r.Context(), author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// refresh authors page
	http.Redirect(w, r, "/authors", http.StatusSeeOther)
}
