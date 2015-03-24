package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"os/exec"
	"path"
)

func flags() (string, string, string, string, error) {
	var database, dsn, packageName, srcPath string

	flag.StringVar(&database, "database", "", "Database name")
	flag.StringVar(&dsn, "dsn", "", "Data Source Name i.e. for Mysql user:pass@protocol(address)")
	flag.StringVar(&packageName, "package", "", "Package name for code generation")
	flag.StringVar(&srcPath, "src", "", "GOPATH relative src path to put the generated code")
	flag.Parse()

	if database == "" || dsn == "" || packageName == "" || srcPath == "" {
		flag.PrintDefaults()
		return "", "", "", "", errors.New("values not set")
	}
	return database, dsn, packageName, srcPath, nil
}

func main() {
	database, dsn, packageName, srcPath, err := flags()
	if err == nil {

		goPath := os.Getenv("GOPATH")
		if goPath == "" {
			log.Fatal("GOPATH is not set")
		}

		connect := dsn + "/" + database
		fullPath := path.Join(goPath, "src", srcPath, packageName)

		if err := os.Mkdir(fullPath, 0755); err != nil && os.IsExist(err) == false {
			log.Fatal(err)
		}

		tables, err := ParseMySQLDB(database, connect)
		if err != nil {
			log.Fatal(err)
		}

		for _, t := range tables {
			file, err := os.Create(fullPath + "/" + t.GoName + ".go")
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			if err := t.Generate(file, packageName); err != nil {
				log.Fatal(err)
			}
		}

		out, err := exec.Command("gofmt", "-w", fullPath).CombinedOutput()
		if err != nil {
			log.Println(err)
		}
		os.Stdout.Write(out)
	}
}
