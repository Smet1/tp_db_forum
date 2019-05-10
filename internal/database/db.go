package database

import (
	"github.com/jackc/pgx"
)

var config = pgx.ConnConfig{
	Host: "localhost",
	Port: 5432,
	//Database: "postgres",
	Database: "docker",
	User:     "docker",
	Password: "docker",
}

var Connection *pgx.ConnPool

func init() {
	Connection, _ = pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig:     config,
			MaxConnections: 50,
		})
}
