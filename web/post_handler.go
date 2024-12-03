package web

import (
	"net/http"
	"super-descuentos/components"
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
