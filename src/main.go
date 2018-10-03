package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

const (
	port = ":8080"
)

func main() {

	// Grab the service name from the environment
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		log.Fatalf("SERVICE_NAME env var must be set")
	}

	// Build up the DB hose and name and connect to it
	dbHost := fmt.Sprintf("%s-db", serviceName)
	dbName := serviceName
	connStr := fmt.Sprintf("postgres://postgres:postgres@%s:5432/%s?sslmode=disable", dbHost, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Return a simple message
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling request for /")
		fmt.Fprintf(w, "Hello from %s", serviceName)
	})

	// Query the time on the db and return (to test connectivity)
	http.HandleFunc("/demo", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling request for /demo")

		t, err := queryDB(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Value from postgres: %s", t)
	})

	log.Printf("Starting web server on port %s", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func queryDB(db *sql.DB) (string, error) {
	rows, err := db.Query("SELECT name FROM demo")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var name string
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			return "nil", err
		}
	}

	err = rows.Err()
	return name, err
}
