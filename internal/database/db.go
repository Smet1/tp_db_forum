package database

import (
	"github.com/jackc/pgx"
)

var config = pgx.ConnConfig{
	Host:     "localhost",
	Port:     5432,
	Database: "docker",
	User:     "docker",
	Password: "docker",
}

var Connection, _ = pgx.NewConnPool(
	pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: 1000,
	})
