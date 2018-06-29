package main

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
	// []byte or time.Time
	"time":      "time.Time",
	"timestamp": "time.Time",
	"date":      "time.Time",
	"datetime":  "time.Time",
}

// map of Go types to Null types
var nullTypes = map[string]string{
	"[]byte":    "[]byte",
	"float64":   "sql.NullFloat64",
	"int64":     "sql.NullInt64",
	"string":    "sql.NullString",
	"time":      "sql.NullTime",
	"timestamp": "sql.NullTime",
	"date":      "sql.NullTime",
	"datetime":  "sql.NullTime",
}

func goType(dt string, nullable bool) (string, []string) {
	imports := []string{}
	gotype := dataTypes[dt]
	if gotype == "" {
		gotype = "[]byte"
	}
	if nullable {
		imports = append(imports, "database/sql")
		gt, ok := nullTypes[gotype]
		if !ok {
			gotype = "[]byte"
		} else {
			gotype = gt
		}
	}
	if gotype == "time.Time" {
		imports = append(imports, "time")
	}
	return gotype, imports
}
