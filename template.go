package main

import "text/template"

// The output template to create a Go struct for a db table
const constTemp = `package {{.PackageName}}

import (
	"database/sql"
	{{printf "// blank import for mysql driver"}}
	_ "github.com/go-sql-driver/mysql"
)

{{printf "// A %s is a direct representation of the database table %s" .GoName .DBName}}
type {{.GoName}} struct {
	{{range $d := .Fields}}{{$d.GoDecl}}
	{{end}}}

{{if .GoKeyName}}
{{printf "// Get will query and scan a database row based on the key"}}
func (sc *{{.GoName}})Get(db *sql.DB, key {{.GoKeyType}}) error {
	query := "SELECT {{.SelectFields}} FROM {{.DBName}} WHERE {{.DBKeyName}} = ?"

	rows, err := db.Query(query, key)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	if rows.Scan({{.ScanString}}) != nil {
		return err
	}
	return nil
}

{{printf "// Set will insert or update a row in table %s" .DBName}}
{{printf "// if %s is blank or zero it will insert otherwise it will update" .GoKeyName}}
{{printf "// if %s is an auto_increment field it will not try and insert the field" .GoKeyName}}
func (sc {{.GoName}})Set(db *sql.DB) (sql.Result, error) {
	if {{.GoKeyCompare}}  {
		sqlstr := "UPDATE {{.UpdateFields}} FROM {{.DBName}} WHERE {{.DBKeyName}} = ?"
		res, err := db.Exec(sqlstr, {{.UpdateValues}}, sc.{{.GoKeyName}})
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	sqlstr := "INSERT INTO {{.DBName}} ({{.InsertFields}}) VALUES ({{.InsertPlaceHolders}})"
	res, err := db.Exec(sqlstr, {{.InsertValues}})
	if err != nil {
		return nil, err
	}
	return res, nil
}
{{else}}
{{printf "// Set will insert a new row into" .DBName}}
func (sc {{.GoName}})Set(db *sql.DB) (sql.Result, error) {
	sqlstr := "INSERT INTO {{.DBName}} ({{.InsertFields}}) VALUES ({{.InsertPlaceHolders}})"
	res, err := db.Exec(sqlstr, {{.InsertValues}})
	if err != nil {
		return nil, err
	}
	return res, nil
}
{{end}}
`

var localTemplate *template.Template

// schemaTemplate returns a static template or creates one and returns it
func schemaTemplate() *template.Template {
	if localTemplate == nil {
		localTemplate = template.Must(template.New("Schema Template").Parse(constTemp))
	}
	return localTemplate
}
