package main

import (
	"fmt"
	"strconv"
	"strings"
)

// map of mysql field types to Go types
var dataTypes = map[string]string{
	"int":        "int64",
	"tinyint":    "int64",
	"smallint":   "int64",
	"mediumint":  "int64",
	"bigint":     "int64", // DANGER: if unsigned and > math.MaxInt64 will return a string
	"float":      "float64",
	"double":     "float64",
	"real":       "float64",
	"decimal":    "[]byte",
	"numeric":    "[]byte",
	"varchar":    "string",
	"bit":        "[]byte",
	"enum":       "string",
	"set":        "string",
	"blob":       "[]byte",
	"tinyblob":   "[]byte",
	"mediumblob": "[]byte",
	"longblob":   "[]byte",
	"text":       "string",
	"tinytext":   "string",
	"mediumtext": "string",
	"longtext":   "string",
	"char":       "string",
	"binary":     "[]byte",
	"varbinary":  "[]byte",
	"year":       "int64",
	"time":       "[]byte",
	"timestamp":  "[]byte",
	"date":       "[]byte", // NOTE: could use time.Time
	"datetime":   "[]byte", // NOTE: could use time.Time
}

// map of Go types to Null types
var nullTypes = map[string]string{
	"[]byte":  "[]byte",
	"float64": "sql.NullFloat64",
	"int64":   "sql.NullInt64",
	"string":  "sql.NullString",
}

// A SQLField represents a field in a database table
type SQLField struct {
	GoDecl        string
	GoVar         string
	GoType        string
	DBField       string
	DBType        string
	DBFullType    string
	Nullable      bool
	Key           string
	PrimaryKey    bool
	AutoIncrement bool
}

// NewSQLField returns a new SQLField based on the meta data for the db field
func NewSQLField(row *DBTblInfo) (*SQLField, error) {
	f := new(SQLField)
	f.DBField = row.columnName
	f.DBType = row.dataType
	f.DBFullType = row.columnType
	f.Nullable = false
	if row.isNullable == "YES" {
		f.Nullable = true
	}
	f.setGoType(row)
	f.setGoVar(row)
	f.GoDecl = fmt.Sprintf("%s %s `db:\"%s\"`", f.GoVar, f.GoType, row.columnName)
	f.Key = row.columnKey
	f.PrimaryKey = false
	f.AutoIncrement = false
	if row.columnKey == "PRI" {
		f.PrimaryKey = true
		if row.extra == "auto_increment" {
			f.AutoIncrement = true
		}
	}
	return f, nil
}

// setGoType determines the appropriate Go type based on the db field type
func (f *SQLField) setGoType(row *DBTblInfo) {
	f.GoType = dataTypes[row.dataType]
	if f.GoType == "" {
		f.GoType = "[]byte"
	}
	if f.Nullable {
		f.GoType = nullTypes[f.GoType]
	}
}

// setGoVar creates a valid Go var name from the db field name
func (f *SQLField) setGoVar(row *DBTblInfo) {
	numToWords := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	name := row.columnName

	// Check the first character in the name. Can't use a number
	i, err := strconv.Atoi(name[:1])
	if err == nil {
		name = strings.Replace(name, name[:1], numToWords[i], 1)
	}

	// remove underscores and upper case the first letter
	var ts string
	s := strings.Split(name, "_")
	for _, sub := range s {
		switch sub {
		case "id":
			ts += "ID"
		case "uid":
			ts += "UID"
		case "ip":
			ts += "IP"
		case "api":
			ts += "API"
		case "url":
			ts += "URL"
		default:
			ts += strings.Title(sub)
		}
	}
	f.GoVar = ts
}
