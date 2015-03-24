package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// ParseMySQLDB connects to a mysql database, gets a list of tables and builds
// a Go struct that can be turned into valid Go code representing a db table
func ParseMySQLDB(dbname, connect string) ([]*SQLTable, error) {
	db, err := sql.Open("mysql", connect)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	dbtables, err := db.Query("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ?", dbname)
	if err != nil {
		return nil, err
	}
	defer dbtables.Close()

	var tname string

	tables := make([]*SQLTable, 0, 50)

	for dbtables.Next() {
		if err := dbtables.Scan(&tname); err != nil {
			return nil, err
		}

		tbl, err := MySQLTableInfo(db, dbname, tname)
		if err != nil {
			return nil, err
		}
		tables = append(tables, tbl)
	}

	if err = dbtables.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

// MySQLTableInfo queries the db for table and field meta data
func MySQLTableInfo(db *sql.DB, dbName, tblName string) (*SQLTable, error) {
	query := "SELECT COLUMN_NAME, DATA_TYPE, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY, EXTRA FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = ? AND TABLE_SCHEMA = ? ORDER BY ORDINAL_POSITION"
	rows, err := db.Query(query, tblName, dbName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fieldInfo := make([]*DBTblInfo, 0, 10)
	for rows.Next() {
		info := &DBTblInfo{}
		err := rows.Scan(&info.columnName, &info.dataType, &info.columnType, &info.isNullable, &info.columnKey, &info.extra)
		if err != nil {
			return nil, err
		}
		fieldInfo = append(fieldInfo, info)
	}

	tbl, err := NewSQLTable(dbName, tblName, fieldInfo)
	if err != nil {
		return nil, err
	}

	return tbl, nil
}
