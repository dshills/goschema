package main

import (
	"bytes"
	"testing"
)

func TestOutput(t *testing.T) {
	conf := config{
		User:        "admin",
		Password:    "abc",
		Host:        "localhost",
		Port:        3306,
		Name:        "redivus",
		PackageName: "model",
	}
	db, err := dbConnect(&conf)
	if err != nil {
		t.Fatal(err)
	}
	tbls, err := getTables(db, conf.Name)
	if err != nil {
		t.Fatal(err)
	}

	var b bytes.Buffer
	for _, tb := range tbls {
		ts, imports := goStruct(tb)
		err := writeModel(&b, imports, ts, conf.PackageName, goName(tb.Name))
		if err != nil {
			t.Error(err)
		}
	}
}

func TestFileOutput(t *testing.T) {
	conf := config{
		User:        "admin",
		Password:    "abc",
		Host:        "localhost",
		Port:        3306,
		Name:        "redivus",
		PackageName: "model",
		OutputDir:   "./_output",
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
		ts, imports := goStruct(tb)
		err := writeFile(conf.OutputDir, imports, ts, conf.PackageName, goName(tb.Name))
		if err != nil {
			t.Error(err)
		}
	}
}
