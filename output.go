package main

import (
	"fmt"
	"io"
	"os"
)

func writeModel(wr io.Writer, imports, gostruct, pkgname, tblname string) error {
	str := fmt.Sprintf("package %s\n\n", pkgname)
	_, err := wr.Write([]byte(str))
	if err != nil {
		return err
	}
	_, err = wr.Write([]byte(imports))
	if err != nil {
		return err
	}
	str = fmt.Sprintf("// %v is database model\n", tblname)
	_, err = wr.Write([]byte(str))
	if err != nil {
		return err
	}
	_, err = wr.Write([]byte(gostruct))
	return err
}

func writeFile(path, imports, gostruct, pkgname, tblname string) error {
	fname := fmt.Sprintf("%v/%v.go", path, tblname)
	file, err := os.Create(fname)
	if err != nil {
		return err
	}
	return writeModel(file, imports, gostruct, pkgname, tblname)
}
