package main

import (
	"encoding/json"
	"os"
)

type config struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Name        string `json:"name"`
	User        string `json:"user"`
	Password    string `json:"password"`
	PackageName string `json:"package_name"`
	OutputDir   string `json:"output_dir"`
	DataTypes   []struct {
		DataType string `json:"data_type"`
		GoType   string `json:"go_type"`
	} `json:"data_types"`
	NullTypes []struct {
		DataType string `json:"data_type"`
		NullType string `json:"null_type"`
	} `json:"null_types"`
}

func readConfig(path string) (*config, error) {
	conf := config{}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&conf)
	return &conf, err
}
