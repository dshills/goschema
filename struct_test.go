package main

import (
	"testing"
)

func TestGoStruct(t *testing.T) {
	opt := options{
		user:     "admin",
		password: "abc",
		host:     "localhost",
		port:     3306,
		name:     "redivus",
	}
	db, err := dbConnect(&opt)
	if err != nil {
		t.Fatal(err)
	}
	tbls, err := getTables(db, opt.name)
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
