package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type options struct {
	user        string
	password    string
	host        string
	port        int
	packageName string
	name        string
}

func flags() (*options, error) {
	opt := options{}
	flag.StringVar(&opt.user, "user", "", "User for login")
	flag.StringVar(&opt.name, "name", "", "Database name")
	flag.StringVar(&opt.password, "password", "", "Password to use when connecting to server")
	flag.StringVar(&opt.host, "host", "127.0.0.1", "Connect to host")
	flag.IntVar(&opt.port, "port", 3306, "Port number to use for connection")
	flag.StringVar(&opt.packageName, "package", "", "Package name for code generation")
	flag.Parse()

	if opt.user == "" || opt.packageName == "" || opt.name == "" {
		flag.Usage()
		return nil, fmt.Errorf("Not set")
	}
	return &opt, nil
}

func main() {
	opt, err := flags()
	if err != nil {
		os.Exit(1)
	}

	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", opt.user, opt.password, opt.host, opt.port, opt.name)

	tables, err := ParseMySQLDB(opt.name, dsn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, t := range tables {
		file, err := os.Create(filepath.Join(path, t.GoName+".go"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()

		if err := t.Generate(file, opt.packageName); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	out, err := exec.Command("gofmt", "-w", path).CombinedOutput()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Stdout.Write(out)
}
