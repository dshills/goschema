package main

import (
	"io"
	"strconv"
	"strings"
)

// A DBTblInfo contains information about a mysql table
type DBTblInfo struct {
	columnName string
	dataType   string
	columnType string
	isNullable string
	columnKey  string
	extra      string
}

// A SQLTable represents a potential Go struct based on a table in a database
type SQLTable struct {
	PackageName   string
	GoName        string
	DBName        string
	GoKeyType     string
	GoKeyName     string
	DBKeyName     string
	AutoIncrement bool
	Fields        []*SQLField
}

// NewSQLTable returns a SQLTable with information about the table in the database
func NewSQLTable(dbName, tblName string, fieldInfo []*DBTblInfo) (*SQLTable, error) {
	tbl := new(SQLTable)
	tbl.DBName = tblName
	tbl.setGoName(tblName)
	tbl.Fields = make([]*SQLField, 0, len(fieldInfo))
	tbl.PackageName = "schema"

	foundKey := false
	for _, row := range fieldInfo {
		fld, err := NewSQLField(row)
		if err != nil {
			return nil, err
		}
		if fld.PrimaryKey && foundKey {
			// no support for multiple primary keys
			tbl.DBKeyName = ""
			tbl.GoKeyName = ""
			tbl.GoKeyType = ""
		} else if fld.PrimaryKey {
			tbl.DBKeyName = fld.DBField
			tbl.GoKeyName = fld.GoVar
			tbl.GoKeyType = fld.GoType
			tbl.AutoIncrement = fld.AutoIncrement
			foundKey = true
		}
		tbl.Fields = append(tbl.Fields, fld)
	}

	return tbl, nil
}

// setGoName creates a valid Go struct name based on a db table name
func (tbl *SQLTable) setGoName(name string) {
	numToWords := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

	i, err := strconv.Atoi(name[:1])
	if err == nil {
		name = strings.Replace(name, name[:1], numToWords[i], 1)
	}

	var ts string
	s := strings.Split(name, "_")
	for _, sub := range s {
		ts += strings.Title(sub)
	}
	tbl.GoName = ts
}

// SelectFields returns a comma seperated list of table fields for use in a select statement
// It is used in the template to build a select statement
func (tbl SQLTable) SelectFields() string {
	flds := make([]string, 0, 10)
	for _, f := range tbl.Fields {
		flds = append(flds, f.DBField)
	}
	return strings.Join(flds, ", ")
}

// ScanString returns a formatted string of fields names passed to the scan function
// it is used in the template to build the query row scan
func (tbl SQLTable) ScanString() string {
	vars := make([]string, 0, 10)
	for _, f := range tbl.Fields {
		vars = append(vars, "&sc."+f.GoVar)
	}
	return strings.Join(vars, ", ")
}

// InsertFields returns a formated string of table fields to use in an insert statement
// It is used in the template to build an insert statement
func (tbl SQLTable) InsertFields() string {
	flds := make([]string, 0, 10)
	for _, f := range tbl.Fields {
		if !f.AutoIncrement {
			flds = append(flds, f.DBField)
		}
	}
	return strings.Join(flds, ", ")
}

// InsertPlaceHolders returns a formated string of ? placeholders to use in an insert statement
// It is used in the template to build an insert statement
func (tbl SQLTable) InsertPlaceHolders() string {
	flds := make([]string, 0, 10)
	for _, f := range tbl.Fields {
		if !f.AutoIncrement {
			flds = append(flds, "?")
		}
	}
	return strings.Join(flds, ", ")
}

// InsertValues returns a formated string of field values to use in an insert statement
// It is used in the template to build an insert statement
func (tbl SQLTable) InsertValues() string {
	vars := make([]string, 0, 10)
	for _, f := range tbl.Fields {
		if !f.AutoIncrement {
			vars = append(vars, "sc."+f.GoVar)
		}
	}
	return strings.Join(vars, ", ")
}

// UpdateFields returns a formated string of fields to use in an update statement
// It is used in the template to build an update statement
func (tbl SQLTable) UpdateFields() string {
	flds := make([]string, 0, 10)
	for _, f := range tbl.Fields {
		if !f.PrimaryKey {
			flds = append(flds, f.DBField+"=?")
		}
	}
	return strings.Join(flds, ", ")
}

// UpdateValues returns a formated string of field values to use in an update statement
// It is used in the template to build an update statement
func (tbl SQLTable) UpdateValues() string {
	flds := make([]string, 0, 10)
	for _, f := range tbl.Fields {
		if !f.PrimaryKey {
			flds = append(flds, "sc."+f.GoVar)
		}
	}
	return strings.Join(flds, ", ")
}

// GoKeyCompare returns a string to compare for an empty value based on field type
// It is used in the template to check for empty values
func (tbl SQLTable) GoKeyCompare() string {
	switch tbl.GoKeyType {
	case "int64", "float64":
		return "sc." + tbl.GoKeyName + "> 0"
	case "sql.NullInt64":
		return "sc." + tbl.GoKeyName + ".Int64 > 0 && " + tbl.GoKeyName + ".Valid"
	case "sql.NullFloat64":
		return "sc." + tbl.GoKeyName + ".Float64 > 0 && " + tbl.GoKeyName + ".Valid"
	case "sql.NullString":
		return "len(sc." + tbl.GoKeyName + ".String) > 0 && " + tbl.GoKeyName + ".Valid"
	}
	return "len(sc." + tbl.GoKeyName + ") > 0"
}

// Generate writes the Go representation of a table to the supplied io.Writer
// It uses the packagename in the template for creating code
func (tbl *SQLTable) Generate(wr io.Writer, PackageName string) error {
	tbl.PackageName = PackageName
	tmpl := schemaTemplate()
	err := tmpl.Execute(wr, tbl)
	if err != nil {
		return err
	}
	return nil
}
