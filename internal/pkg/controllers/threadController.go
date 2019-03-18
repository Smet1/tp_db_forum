package controllers

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
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
	log.Println("\t", t)


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

	log.Println(createdThread)

	ResponseObject(res, http.StatusCreated, createdThread)
}

func GetThreads(res http.ResponseWriter, req *http.Request) {
	slugName, err := checkVar("slug", req)
	if err != nil {
		ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get user slug").Error())
		return
	}

	query := req.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	since := query.Get("since")
	desc,_ := strconv.ParseBool(query.Get("desc"))

	fmt.Println(query)
	fmt.Println(limit)
	fmt.Println(since)
	fmt.Println(desc)

	threads, err, status := models.GetForumThreads(slugName.(string), limit, since, desc)
	if err != nil {
		if status == http.StatusNotFound {
			ErrResponse(res, status, err.Error())
			return
		}

		ErrResponse(res, status, err.Error())
		return
	}

	ResponseObject(res, http.StatusOK, threads)
}
