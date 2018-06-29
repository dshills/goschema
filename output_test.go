package main

import (
	"bytes"
	"testing"
)

func TestOutput(t *testing.T) {
	opt := options{
		user:        "admin",
		password:    "abc",
		host:        "localhost",
		port:        3306,
		name:        "redivus",
		packageName: "model",
	}
	db, err := dbConnect(&opt)
	if err != nil {
		t.Fatal(err)
	}
	tbls, err := getTables(db, opt.name)
	if err != nil {
		t.Fatal(err)
	}

	var b bytes.Buffer
	for _, tb := range tbls {
		ts, imports := goStruct(tb)
		err := writeModel(&b, imports, ts, opt.packageName, goName(tb.Name))
		if err != nil {
			t.Error(err)
		}
	}
}

func TestFileOutput(t *testing.T) {
	opt := options{
		user:        "admin",
		password:    "abc",
		host:        "localhost",
		port:        3306,
		name:        "redivus",
		packageName: "model",
		output:      "./_output",
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
		ts, imports := goStruct(tb)
		err := writeFile(opt.output, imports, ts, opt.packageName, goName(tb.Name))
		if err != nil {
			t.Error(err)
		}
	}
}
