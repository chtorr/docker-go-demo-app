package main

import (
	"log"
	"testing"
)

func TestIntegration_QueryDB(t *testing.T) {
	serviceName, err := getServiceName()
	if err != nil {
		t.Fatal(err)
		log.Fatal(err)
	}

	db, err := getDB(serviceName)
	if err != nil {
		t.Fatal(err)
	}

	name, err := queryDB(db)
	if err != nil {
		t.Fatal(err)
	}

	expected := "test name"
	if name != expected {
		t.Fatalf("Expected (%s) found (%s)", expected, name)
	}

}
