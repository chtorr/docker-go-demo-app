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

func getServiceName() (string, error) {
	// Grab the service name from the environment
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		return "", fmt.Errorf("SERVICE_NAME env var must be set")
	}
	return serviceName, nil
}

func main() {

	serviceName, err := getServiceName()
	if err != nil {
		log.Fatal(err)
	}

	db, err := getDB(serviceName)
	if err != nil {
		log.Fatal(err)
	}

	// Return a simple message
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling request for /")
		fmt.Fprintf(w, "Hello from %s", serviceName)
	})

	// Query some data from the db and return (to test connectivity and migrations)
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

func getDB(serviceName string) (*sql.DB, error) {
	// Build up the DB hose and name and connect to it
	dbHost := fmt.Sprintf("%s-db", serviceName)
	dbName := serviceName
	connStr := fmt.Sprintf("postgres://postgres:postgres@%s:5432/%s?sslmode=disable", dbHost, dbName)
	return sql.Open("postgres", connStr)
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
