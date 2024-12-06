package web

import (
	"net/http"
	"super-descuentos/components"
	"super-descuentos/errs"
	"super-descuentos/model"
	"super-descuentos/utils"

	"github.com/google/uuid"
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
