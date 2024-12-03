package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"super-descuentos/relational"
	"super-descuentos/rest"
	"super-descuentos/web"

	_ "modernc.org/sqlite"
)

func main() {
	// sqlite db
	db, err := sql.Open("sqlite", "./super-descuentos.db")
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
	hostName := "0.0.0.0:8080"
	fmt.Printf("Starting server on %s...\n", hostName)
	if err := http.ListenAndServe(hostName, handler); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
