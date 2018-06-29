package main

import (
	"fmt"
)

func goStruct(table DBTable) (string, string) {
	structStr := fmt.Sprintf("type %s struct {\n", goName(table.Name))
	imap := make(map[string]bool)
	for _, fld := range table.Fields {
		f := goName(fld.ColumnName)
		nullable := false
		if fld.IsNullable == "YES" {
			nullable = true
		}
		ty, imp := goType(fld.DataType, nullable)
		for _, im := range imp {
			imap[im] = true
		}
		structStr += fmt.Sprintf("\t%s %s `db:\"%s\"`\n", f, ty, fld.ColumnName)
	}
	structStr += "}\n"

	imports := ""
	if len(imap) > 0 {
		imports = "import (\n"
		for k := range imap {
			imports += "\t\"" + k + "\"\n"
		}
		imports += ")\n"
	}
	return structStr, imports
}
