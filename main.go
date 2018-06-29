package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	// blank import for mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type options struct {
	name        string
	user        string
	password    string
	host        string
	port        int
	packageName string
	output      string
}

func flags() (*options, error) {
	opt := options{}
	flag.StringVar(&opt.user, "user", "root", "User for login")
	flag.StringVar(&opt.name, "name", "", "Database name")
	flag.StringVar(&opt.password, "password", "", "Password to use when connecting to server")
	flag.StringVar(&opt.host, "host", "127.0.0.1", "Connect to host")
	flag.IntVar(&opt.port, "port", 3306, "Port number to use for connection")
	flag.StringVar(&opt.packageName, "package", "main", "Package name for code generation")
	flag.StringVar(&opt.output, "output", "./", "Output directory")
	flag.Parse()

	if opt.name == "" {
		flag.Usage()
		return nil, fmt.Errorf("Not set")
	}
	return &opt, nil
}

func main() {
	opt, err := flags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db, err := dbConnect(opt)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tables, err := getTables(db, opt.name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, tbl := range tables {
		ts, imports := goStruct(tbl)
		err := writeFile(opt.output, imports, ts, opt.packageName, goName(tbl.Name))
		if err != nil {
			fmt.Println(err)
		}
	}

	if err = exec.Command("gofmt", "-w", opt.output).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func dbConnect(opt *options) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", opt.user, opt.password, opt.host, opt.port, opt.name)
	return sqlx.Connect("mysql", dsn)
}
