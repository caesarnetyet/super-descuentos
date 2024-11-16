package main

import (
	"fmt"
	"net/http"
	"os"
	"super-descuentos/rest"
	"super-descuentos/store"
)

func main() {
	store := store.NewInMemoryStore()
	server := rest.NewServer(store)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Disable automatic HTTPS redirect in test environment
	if os.Getenv("TEST_ENV") == "true" {
		http.DefaultServeMux = http.NewServeMux()
	}

	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Servidor corriendo en http://localhost%s\n", addr)

	err := http.ListenAndServe(addr, server)
	if err != nil {
		fmt.Println(err)
	}
}
