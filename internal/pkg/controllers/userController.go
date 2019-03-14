package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"tp_db_forum/internal/pkg/models"
)

type getUserProfileResponse struct {
	Message string `json:"message"`
}

func GetUserProfile(ctx *fasthttp.RequestCtx) {
	name := ctx.UserValue("nickname")
	_, _ = fmt.Fprint(ctx, name)

	u, err := models.GetUserByNickname(name.(string))
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(404)
		msg := make([]byte, 0, 1)
		json.Unmarshal(msg, getUserProfileResponse{"Can't find user with nickname " + name.(string)})
		ctx.Write(msg)
	}
}

func UpdateUserProfile(ctx *fasthttp.RequestCtx) {

}

func CreateUser(ctx *fasthttp.RequestCtx) {
	fmt.Println("POST --- create")
}
