package controllers

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
	"tp_db_forum/internal/pkg/models"
)

func PrintThread(t models.Thread) {
	fmt.Println("\tauthor = ", t.Author)
	fmt.Println("\tcreated = ", t.Created)
	fmt.Println("\tforum = ", t.Forum)
	fmt.Println("\tid = ", t.ID)
	fmt.Println("\tmessage = ", t.Message)
	fmt.Println("\tslug = ", t.Slug)
	fmt.Println("\ttitle = ", t.Title)
	fmt.Println("\tvotes = ", t.Votes)
}

func CreateThread(res http.ResponseWriter, req *http.Request) {
	log.Println("=============")
	log.Println("CreateThread", req.URL)

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
	fmt.Println("\tGET")
	PrintThread(t)

	createdThread, err, status := models.CreateThread(t)
	if err != nil {
		if status == http.StatusNotFound {
			ErrResponse(res, status, err.Error())

			log.Println(err.Error())
			return
		}

		if status == http.StatusInternalServerError {
			ErrResponse(res, status, err.Error())

			log.Println(err.Error())
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

	fmt.Println("\tCREATED")
	PrintThread(createdThread)

	ResponseObject(res, http.StatusCreated, createdThread)
}

func GetThreads(res http.ResponseWriter, req *http.Request) {
	log.Println("=============")
	log.Println("GetThreads", req.URL)

	slugName, err := checkVar("slug", req)
	if err != nil {
		ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get user slug").Error())
		return
	}

	query := req.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	since := query.Get("since")
	desc, _ := strconv.ParseBool(query.Get("desc"))

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
	fmt.Println("OK RESP")
	for i, val := range threads {
		fmt.Println("--", i, "--")
		PrintThread(val)
	}
}
