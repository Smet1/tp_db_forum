package controllers

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
	"tp_db_forum/internal/pkg/models"
)

func CreatePosts(res http.ResponseWriter, req *http.Request) {
	log.Println("=============")
	log.Println("CreatePosts", req.URL)

	slugOrId, err := checkVar("slug_or_id", req)
	if err != nil {
		ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get user slug").Error())
		return
	}
	slug := slugOrId.(string)
	id, err := strconv.ParseInt(slug, 10, 32)
	if id == 0 {
		id = -1
	}

	existingThread, err, status := models.GetThreadByIDorSlug(int(id), slug)
	if err != nil {
		ErrResponse(res, status, errors.Wrap(err, "not found").Error())
		return
	}

	fmt.Println(existingThread)

	postsToCreate := make([]models.Post, 0, 1)

	status, err = ParseRequestIntoStruct(req, &postsToCreate)
	if err != nil {
		ErrResponse(res, status, err.Error())

		log.Println("\t", errors.Wrap(err, "ParseRequestIntoStruct error"))
		return
	}

	createdPosts, err, status := models.CreatePosts(postsToCreate, existingThread)
	if err != nil {
		ErrResponse(res, status, err.Error())

		log.Println("\t", errors.Wrap(err, "models.CreatePost error"))
		return
	}

	ResponseObject(res, http.StatusCreated, createdPosts)
}

func GetThreadPosts(res http.ResponseWriter, req *http.Request) {
	log.Println("=============")
	log.Println("GetThreadPosts", req.URL)

	slugOrId, err := checkVar("slug_or_id", req)
	if err != nil {
		ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get user slug").Error())
		return
	}
	slug := slugOrId.(string)
	id, err := strconv.ParseInt(slug, 10, 32)
	if id == 0 {
		id = -1
	}

	existingThread, err, status := models.GetThreadByIDorSlug(int(id), slug)
	if err != nil {
		ErrResponse(res, status, errors.Wrap(err, "not found").Error())
		return
	}

	fmt.Println(existingThread)

	query := req.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	since, _ := strconv.Atoi(query.Get("since"))
	sort := query.Get("sort")
	desc, _ := strconv.ParseBool(query.Get("desc"))

	fmt.Println(query)
	fmt.Println(limit)
	fmt.Println(since)
	fmt.Println(sort)
	fmt.Println(desc)

	sortedPosts, err, status := models.GetSortedPosts(existingThread, limit, since, sort, desc)
	if err != nil {
		ErrResponse(res, status, err.Error())

		return
	}

	ResponseObject(res, status, sortedPosts)
}

func UpdatePost(res http.ResponseWriter, req *http.Request) {
	log.Println("=============")
	log.Println("UpdatePost", req.URL)

	slug, err := checkVar("id", req)
	if err != nil {
		ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get user slug").Error())
		return
	}
	id, err := strconv.ParseInt(slug.(string), 10, 64)
	if id == 0 {
		id = -1
	}

	existingPost, err, status := models.GetPostByID(id)
	if err != nil {
		ErrResponse(res, status, err.Error())

		return
	}

	newPost := models.Post{}

	status, err = ParseRequestIntoStruct(req, &newPost)
	if err != nil {
		ErrResponse(res, status, err.Error())

		log.Println("\t", errors.Wrap(err, "ParseRequestIntoStruct error"))
		return
	}

	updatedPost, err, status := models.UpdatePost(existingPost, newPost)

	if err != nil {
		ErrResponse(res, status, err.Error())

		return
	}

	ResponseObject(res, status, updatedPost)
}
