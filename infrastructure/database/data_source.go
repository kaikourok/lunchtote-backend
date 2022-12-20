package database

import (
	"fmt"
)

type DataSource struct {
	Host     string
	Port     int
	User     string
	Database string
	Password string
}

type dataSourceOptionStruct struct {
	sslmode string
}

type dataSourceOption func(*dataSourceOptionStruct)

func dsnOptionSslMode(mode string) dataSourceOption {
	return func(option *dataSourceOptionStruct) {
		option.sslmode = mode
	}
}

func (d *DataSource) getDataSourceString(options ...dataSourceOption) string {
	option := dataSourceOptionStruct{
		sslmode: "disable",
	}

	for _, f := range options {
		f(&option)
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		d.Host,
		d.Port,
		d.User,
		d.Database,
		d.Password,
		option.sslmode,
	)
}
