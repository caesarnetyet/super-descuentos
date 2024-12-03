package web

import (
	"net/http"
	"super-descuentos/components"
	"super-descuentos/utils"
)

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
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

	component := components.Layout("Home", components.HomePage(posts))
	component.Render(r.Context(), w)

}
