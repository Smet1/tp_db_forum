package controllers

import (
	"github.com/pkg/errors"
	"log"
	"net/http"
	"tp_db_forum/internal/pkg/models"
)

func CreateThread(res http.ResponseWriter, req *http.Request) {
	slugName, err := checkVar("slug", req)
	if err != nil {
		ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get user slug").Error())
		return
	}
	t := models.Thread{}

	status, err := ParseRequestIntoStruct(req, &t)
	if err != nil {
		ErrResponse(res, status, err.Error())

		log.Println("\t", errors.Wrap(err, "ParseRequestIntoStruct error"))
		return
	}
	t.Forum = slugName.(string)

	createdThread, err, status := models.CreateThread(t)
	if err != nil {
		if status == http.StatusNotFound {
			ErrResponse(res, status, err.Error())
			return
		}

		//existingForum, err := models.GetForumBySlug(f.Slug)
		//if err != nil {
		//	ErrResponse(res, status, err.Error())
		//	return
		//}
		//
		//ResponseObject(res, http.StatusConflict, existingForum)
		//return
	}

	ResponseObject(res, http.StatusCreated, createdThread)

}
