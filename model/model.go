package model

import (
	"errors"
	"strings"

	"github.com/tiamxu/kit/sql"
)

var postgresHandler = sql.NewPreDB()
var mysqlHandler = sql.NewPreDB()
var clickHouseHandler = sql.NewPreDB()

func Init(cfg *sql.Config) error {
	var err error
	switch strings.ToLower(cfg.Driver) {
	case "mysql":
		err = mysqlHandler.Init(cfg)
	case "postgres":
		err = postgresHandler.Init(cfg)
	case "clickhouse":
		err = clickHouseHandler.Init(cfg)
	default:
		return errors.New("unknown driver")
	}

	if err != nil {
		return err
	}

	return nil
}

func GetMysqlDB() *sql.DB {
	return mysqlHandler.DB
}

func GetPostgresDB() *sql.DB {
	return postgresHandler.DB
}

func GetClickhouseDB() *sql.DB {
	return clickHouseHandler.DB
}

var (
	AlertRecordsTableName = "\"alert_records\""
)
