package web

import (
	"net/http"
	"super-descuentos/components"
	"super-descuentos/errs"
	"super-descuentos/model"
)

func (s *Server) handleAuthorForm(w http.ResponseWriter, r *http.Request) {
	component := components.Layout("Authors", components.AuthorsPage())
	err := component.Render(r.Context(), w)

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
	var errs model.ValidationErrors

	if r.Name == "" {
		errs = append(errs, model.ValidationError{
			Field:   "name",
			Message: "el nombre es requerido",
		})
	}

	if r.Email == "" {
		errs = append(errs, model.ValidationError{
			Field:   "email",
			Message: "el correo electrónico es requerido",
		})
	}

	return errs
}

func (s *Server) handleCreateAuthorForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, errs.ErrInvalidFormData.Error(), http.StatusBadRequest)
		return
	}

	author := model.User{
		Name:  r.FormValue("name"),
		Email: r.FormValue("email"),
	}

	if errs := author.Validate(); len(errs) > 0 {
		http.Error(w, errs.Error(), http.StatusBadRequest)
		return
	}

	err = s.store.CreateAuthor(r.Context(), author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
