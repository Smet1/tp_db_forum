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

	router.HandleFunc("/user/{nickname}/profile", controllers.GetUserProfile).Methods("GET")
	router.HandleFunc("/user/{nickname}/profile", controllers.UpdateUserProfile).Methods("POST")
	router.HandleFunc("/user/{nickname}/create", controllers.CreateUser).Methods("POST")

	err := http.ListenAndServe(address, router)
	if err != nil {
		return errors.Wrap(err, "server Run error")
	}

	//router := fasthttprouter.New()
	//
	//router.GET("/user/:nickname/profile", controllers.GetUserProfile)
	//router.POST("/user/:nickname/profile", controllers.GetUserProfile)
	//router.POST("/user/:nickname/create", controllers.CreateUser)
	//
	//err := fasthttp.ListenAndServe(address, router.Handler)
	//if err != nil {
	//	return errors.Wrap(err, "server Run error")
	//}

	return nil
}
