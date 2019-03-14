package main

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"tp_db_forum/internal/app/server"
)

func main() {
	fmt.Println("server will start on port 5000")

	err := server.Run("5000")
	if err != nil {
		log.Fatal(errors.Wrap(err, "cant start server"))
	}
}
