package rest

import (
	"net/http"
	"super-descuentos/model"
	"super-descuentos/utils"
)

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/web/posts", http.StatusFound)
}

func (server *Server) handlePostsPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	offset, limit, err := utils.GetOffsetAndLimit(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	posts, err := server.store.GetPosts(r.Context(), offset, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		Posts []model.Post
	}{
		Posts: posts,
	}
	err = server.templates.ExecuteTemplate(w, "posts/posts.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
