package main

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
	"tp_db_forum/db"
	"tp_db_forum/models"
)



func GetUserProfile(ctx *fasthttp.RequestCtx)  {
	name := ctx.UserValue("nickname")
	_, _ = fmt.Fprint(ctx, name)

	user := models.User{}
	conn := db.Connection
	res, _ := conn.Query(`SELECT about, email, fullname, nickname FROM forum_users
	WHERE nickname = $1`, name)

	if res.Next() {
		err := res.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
		fmt.Println(err)
	}
	fmt.Println(user)

}

func UpdateUserProfile(ctx *fasthttp.RequestCtx)  {

}

func CreateUser(ctx *fasthttp.RequestCtx)  {
	fmt.Println("POST --- create")
}

func main() {
	router := fasthttprouter.New()
	router.GET("/user/:nickname/profile", GetUserProfile)
	router.POST("/user/:nickname/profile", UpdateUserProfile)
	router.POST("/user/:nickname/create", CreateUser)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}


