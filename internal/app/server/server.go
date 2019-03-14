package server

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"tp_db_forum/internal/pkg/controllers"
)

func Run(port string) error {
	address := ":" + port

	router := fasthttprouter.New()

	router.GET("/user/:nickname/profile", controllers.GetUserProfile)
	router.POST("/user/:nickname/profile", controllers.GetUserProfile)
	router.POST("/user/:nickname/create", controllers.CreateUser)

	err := fasthttp.ListenAndServe(address, router.Handler)
	if err != nil {
		return errors.Wrap(err, "server Run error")
	}

	return nil
}
