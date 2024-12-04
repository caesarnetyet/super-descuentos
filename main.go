package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"super-descuentos/relational"
	"super-descuentos/rest"
	"super-descuentos/web"

	_ "modernc.org/sqlite"
)

func main() {
	dbFilePath := "./super-descuentos.db"

	// sqlite db
	db, err := initializeDatabase(dbFilePath)
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
	fmt.Printf("Starting server on http://%s...\n", hostName)
	if err := http.ListenAndServe(hostName, handler); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func initializeDatabase(dbFilePath string) (*sql.DB, error) {
	// Check if the database file exists
	_, err := os.Stat(dbFilePath)
	dbExists := !os.IsNotExist(err)

	// Open the SQLite database
	db, err := sql.Open("sqlite", dbFilePath)
	if err != nil {
		return nil, err
	}

	// If the database file doesn't exist, create it and run schema and seed SQL files
	if !dbExists {
		fmt.Println("Database file not found. Creating new database and running schema and seed SQL files...")

		// Read and execute schema SQL file
		schemaSQL, err := os.ReadFile("./relational/schema/schema.sql")
		if err != nil {
			return nil, fmt.Errorf("error reading schema.sql: %w", err)
		}
		if _, err := db.Exec(string(schemaSQL)); err != nil {
			return nil, fmt.Errorf("error executing schema.sql: %w", err)
		}

		// Read and execute seed SQL file
		seedSQL, err := os.ReadFile("./relational/seed/seed.sql")
		if err != nil {
			return nil, fmt.Errorf("error reading seed.sql: %w", err)
		}
		if _, err := db.Exec(string(seedSQL)); err != nil {
			return nil, fmt.Errorf("error executing seed.sql: %w", err)
		}
	}

	return db, nil
}
