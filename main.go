package main

import (
	"flag"
	"fmt"
	"os"

	// blank import for mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func flags() (*config, error) {
	confPath := ""
	conf := config{}
	flag.StringVar(&confPath, "config", "", "Configuration file path")
	flag.StringVar(&conf.User, "user", "root", "User for login")
	flag.StringVar(&conf.Name, "name", "", "Database name")
	flag.StringVar(&conf.Password, "password", "", "Password to use when connecting to server")
	flag.StringVar(&conf.Host, "host", "127.0.0.1", "Connect to host")
	flag.IntVar(&conf.Port, "port", 3306, "Port number to use for connection")
	flag.StringVar(&conf.PackageName, "package", "main", "Package name for code generation")
	flag.StringVar(&conf.OutputDir, "output", "./", "Output directory")
	flag.Parse()

	if confPath != "" {
		return readConfig(confPath)
	}

	if conf.Name == "" {
		flag.Usage()
		return nil, fmt.Errorf("Not set")
	}
	return &conf, nil
}

func main() {
	conf, err := flags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db, err := dbConnect(conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, nt := range conf.NullTypes {
		nullTypes[nt.DataType] = nt.NullType
	}
	for _, dt := range conf.DataTypes {
		dataTypes[dt.DataType] = dt.GoType
	}

	tables, err := getTables(db, conf.Name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, tbl := range tables {
		ts, imports := goStruct(tbl)
		err := writeFile(conf.OutputDir, imports, ts, conf.PackageName, goName(tbl.Name))
		if err != nil {
			fmt.Println(err)
		}
	}
	/*
		if err = exec.Command("gofmt", "-w", conf.OutputDir).Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	*/
}

func dbConnect(conf *config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", conf.User, conf.Password, conf.Host, conf.Port, conf.Name)
	return sqlx.Connect("mysql", dsn)
}
