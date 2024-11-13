package main

import (
	"fmt"
	"net/http"
	"super-descuentos/rest"
	"super-descuentos/store"
)

func main() {
	store := store.NewInMemoryStore()
	server := rest.NewServer(store)

	fmt.Println("Servidor corriendo en http://localhost:8080")
	err := http.ListenAndServe(":8080", server)
	if err != nil {
		fmt.Println(err)
	}
}
