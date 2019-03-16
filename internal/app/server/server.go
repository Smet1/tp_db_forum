package server

import (
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"tp_db_forum/internal/pkg/controllers"
)

func Run(port string) error {
	address := ":" + port

	router := mux.NewRouter()

	router.HandleFunc("/api/user/{nickname}/profile", controllers.GetUserProfile).Methods("GET")
	router.HandleFunc("/api/user/{nickname}/profile", controllers.UpdateUserProfile).Methods("POST")
	router.HandleFunc("/api/user/{nickname}/create", controllers.CreateUser).Methods("POST")

	router.HandleFunc("/api/forum/create", controllers.CreateForum).Methods("POST")

	err := http.ListenAndServe(address, router)
	if err != nil {
		return errors.Wrap(err, "server Run error")
	}

	return nil
}
