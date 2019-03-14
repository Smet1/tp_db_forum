package main

import (
	"github.com/pkg/errors"
	"log"
	"tp_db_forum/internal/app/server"
)

func main() {
	err := server.Run("5000")
	if err != nil {
		log.Fatal(errors.Wrap(err, "cant start server"))
	}
}
