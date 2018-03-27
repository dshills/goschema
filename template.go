package main

import "text/template"

func loadTemplate(path string, usesqlx bool) (*template.Template, error) {
	tpath := path + "sql.tmpl"
	if usesqlx {
		tpath = path + "sqlx.tmpl"
	}
	return template.ParseFiles(tpath)
}
