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
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello World"))

	})
	router.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Posts"))
	})
	router.HandleFunc("/author", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Author"))
	})

	server.Handler = router

	return server
}
