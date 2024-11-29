package web

import (
	"net/http"
	"super-descuentos/store"
)

type Server struct {
	store store.Store
	http.Handler
}

func NewServer(store store.Store) *Server {
	server := new(Server)
	server.store = store

	router := http.NewServeMux()
	router.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Posts"))
	})
	router.HandleFunc("/authors", server.handleAuthorForm)
	router.HandleFunc("/", server.handleHome)

	server.Handler = router

	return server
}
