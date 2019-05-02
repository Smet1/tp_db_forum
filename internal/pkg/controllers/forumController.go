package controllers

import (
	"github.com/pkg/errors"
	"net/http"
	"tp_db_forum/internal/pkg/models"
)

func CreateForum(res http.ResponseWriter, req *http.Request) {
	//log.Println("=============")
	//log.Println("CreateForum", req.URL)

	f := models.Forum{}
	status, err := ParseRequestIntoStruct(req, &f)
	if err != nil {
		ErrResponse(res, status, err.Error())

		//log.Println("\t", errors.Wrap(err, "ParseRequestIntoStruct error"))
		return
	}

	createdForum, err := models.CreateForum(f)
	if err != nil {
		_, err := models.GetUserByNickname(f.User)
		if err != nil {
			ErrResponse(res, http.StatusNotFound, "Can't find user with nickname "+f.User)
			return
		}

		//fmt.Println(user)

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

func GetForum(res http.ResponseWriter, req *http.Request) {
	//log.Println("=============")
	//log.Println("GetForum", req.URL)

	searchingSlug, err := checkVar("slug", req)
	if err != nil {
		ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get forum slug").Error())
		return
	}

	f, err := models.GetForumBySlug(searchingSlug.(string))
	if err != nil || f.User == "" {
		ErrResponse(res, http.StatusNotFound, "Can't find slug")
		return
	}

	ResponseObject(res, http.StatusOK, f)
}
