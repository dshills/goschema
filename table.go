package main

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// DBTable represents a table in a database
type DBTable struct {
	Name   string
	Fields []DBField
}

// DBField represents a field in a table
type DBField struct {
	ColumnName string `db:"COLUMN_NAME"`
	DataType   string `db:"DATA_TYPE"`
	ColumnType string `db:"COLUMN_TYPE"`
	IsNullable string `db:"IS_NULLABLE"`
	ColumnKey  string `db:"COLUMN_KEY"`
	Extra      string `db:"EXTRA"`
}

func getTables(db *sqlx.DB, name string) ([]DBTable, error) {
	tblNames, err := tableList(db, name)
	if err != nil {
		return nil, err
	}
	elist := []string{}
	tables := []DBTable{}
	for _, tbl := range tblNames {
		fields, err := tableFields(db, name, tbl)
		if err != nil {
			elist = append(elist, err.Error())
			continue
		}
		table := DBTable{Name: tbl, Fields: fields}
		tables = append(tables, table)
	}
	if len(elist) > 0 {
		return tables, fmt.Errorf(strings.Join(elist, ", "))
	}
	return tables, nil
}

func tableList(db *sqlx.DB, dbname string) ([]string, error) {
	str := "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ?"
	tnames := []string{}
	err := db.Select(&tnames, str, dbname)
	return tnames, err
}

func tableFields(db *sqlx.DB, dbName, tblName string) ([]DBField, error) {
	str := `
	SELECT COLUMN_NAME, DATA_TYPE, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY, EXTRA
	FROM INFORMATION_SCHEMA.COLUMNS
	WHERE TABLE_NAME = ?
		AND TABLE_SCHEMA = ?
	ORDER BY ORDINAL_POSITION`

	fields := []DBField{}
	err := db.Select(&fields, str, tblName, dbName)
	return fields, err
}
