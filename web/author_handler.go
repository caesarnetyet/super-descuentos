package web

import (
	"net/http"
	"super-descuentos/components"
)

func (s *Server) handleAuthorForm(w http.ResponseWriter, r *http.Request) {
	component := components.Layout("Authors", components.AuthorsPage())
	err := component.Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
