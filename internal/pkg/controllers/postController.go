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
