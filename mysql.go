package main

import (
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// ParseMySQLDB connects to a mysql database, gets a list of tables and builds
// a Go struct that can be turned into valid Go code representing a db table
func ParseMySQLDB(db *sqlx.DB, dbname string) ([]SQLTable, error) {
	str := `
	SELECT TABLE_NAME
	FROM INFORMATION_SCHEMA.TABLES
	WHERE TABLE_SCHEMA = ?
	`

	tnames := []string{}
	db.Select(&tnames, str, dbname)

	tables := []SQLTable{}
	elist := []string{}
	for _, t := range tnames {
		tbl, err := MySQLTableInfo(db, dbname, t)
		if err != nil {
			elist = append(elist, err.Error())
			continue
		}
		tables = append(tables, *tbl)
	}

	if len(elist) > 0 {
		return tables, fmt.Errorf(strings.Join(elist, ", "))
	}
	return tables, nil
}

// MySQLTableInfo queries the db for table and field meta data
func MySQLTableInfo(db *sqlx.DB, dbName, tblName string) (*SQLTable, error) {
	str := `
	SELECT COLUMN_NAME, DATA_TYPE, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY, EXTRA
	FROM INFORMATION_SCHEMA.COLUMNS
	WHERE TABLE_NAME = ?
		AND TABLE_SCHEMA = ?
	ORDER BY ORDINAL_POSITION`

	tinfo := []DBTblInfo{}
	err := db.Select(&tinfo, str, tblName, dbName)
	if err != nil {
		return nil, err
	}

	return NewSQLTable(dbName, tblName, tinfo)
}
