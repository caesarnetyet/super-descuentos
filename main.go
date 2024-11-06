package main

import (
	"fmt"
	"net/http"
)

func main() {
	store := NewInMemoryStore()
	server := NewServer(store)

	fmt.Println("Servidor corriendo en http://localhost:8080")
	err := http.ListenAndServe(":8080", server)
	if err != nil {
		fmt.Println(err)
	}
}
