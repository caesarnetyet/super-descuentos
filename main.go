package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"super-descuentos/relational"
	"super-descuentos/rest"
	"super-descuentos/web"
)

func main() {
	// sqlite db
	db, err := sql.Open("sqlite3", "./super-descuentos.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	store := relational.NewSQLStore(db)
	apiServer := rest.NewServer(store)
	webServer := web.NewServer(store)

	// Create the handler for routing
	handler := http.NewServeMux()

	// API handler with prefix stripping
	handler.Handle("/api/", http.StripPrefix("/api", apiServer))

	// Web server handler for all other routes
	handler.Handle("/", webServer)

	// Start the HTTP server
	port := ":8080"
	fmt.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
