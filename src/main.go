package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

const (
	port          = ":8080"
	proxyHostName = "proxy"
	proxyPort     = "80"
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

	// Return a simple message from another HTTP service
	http.HandleFunc("/remote", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling request for /remote")

		remote := r.URL.Query().Get("service")
		if remote == "" {
			http.Error(w, "service name must be provided", http.StatusBadRequest)
			return
		}

		resp, err := queryRemote(remote)
		if err != nil {
			http.Error(w, fmt.Sprintf("error requesting from remote: %v", err), http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Msg from remote (%s): %s", remote, resp)
	})

	log.Printf("Starting web server on port %s", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func queryRemote(serviceName string) (string, error) {
	client := http.Client{}

	url := fmt.Sprintf("http://%s:%s", proxyHostName, proxyPort)
	req, err := http.NewRequest("GET", url, nil)
	req.Host = serviceName
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil

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
