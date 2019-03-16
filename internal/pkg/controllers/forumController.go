package controllers

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"tp_db_forum/internal/pkg/models"
)

func CreateForum(res http.ResponseWriter, req *http.Request) {
	f := models.Forum{}
	status, err := ParseRequestIntoStruct(req, &f)
	if err != nil {
		ErrResponse(res, status, err.Error())

		log.Println("\t", errors.Wrap(err, "ParseRequestIntoStruct error"))
		return
	}

	createdForum, err := models.CreateForum(f)
	if err != nil {
		user, err := models.GetUserByNickname(f.User)
		if err != nil {
			ErrResponse(res, http.StatusNotFound, "Can't find user with nickname " + f.User)
			return
		}

		fmt.Println(user)

		existingForum, err := models.GetForumBySlug(f.Slug)
		if err != nil {
			ErrResponse(res, status, err.Error())
			return
		}

		ResponseObject(res, http.StatusConflict, existingForum)
		return
	}

	ResponseObject(res, http.StatusCreated, createdForum)
}
