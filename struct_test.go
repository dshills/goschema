package main

import (
	"testing"
)

func TestGoStruct(t *testing.T) {
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
	for _, tb := range tbls {
		ts, _ := goStruct(tb)
		if ts == "" {
			t.Errorf("Expected struct got nothing")
		}
	}
}
