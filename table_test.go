package main

import (
	"testing"
)

func TestGetTables(t *testing.T) {
	conf := config{
		User:     "admin",
		Password: "abc",
		Host:     "localhost",
		Port:     3306,
		Name:     "redivus",
	}
	db, err := dbConnect(&conf)
	if err != nil {
		t.Fatal(err)
	}
	tbls, err := getTables(db, conf.Name)
	if err != nil {
		t.Fatal(err)
	}
	if len(tbls) < 10 {
		t.Errorf("Expected > 10 tables got %v\n", len(tbls))
	}
}
