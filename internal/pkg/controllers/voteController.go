package controllers

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
	"tp_db_forum/internal/pkg/models"
)

func CreateVote(res http.ResponseWriter, req *http.Request) {
	log.Println("=============")
	log.Println("CreateVote", req.URL)

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
		ErrResponse(res, status, errors.Wrap(err, "slug not found").Error())

		return
	}

	fmt.Println("--before vote--")
	PrintThread(existingThread)

	voteToCreate := models.Vote{}
	status, err = ParseRequestIntoStruct(req, &voteToCreate)
	if err != nil {
		ErrResponse(res, status, err.Error())

		log.Println("\t", errors.Wrap(err, "ParseRequestIntoStruct error"))
		return
	}

	existingUser, err := models.GetUserByNickname(voteToCreate.Nickname)
	if err != nil {
		ErrResponse(res, http.StatusNotFound, errors.Wrap(err, "user not found").Error())

		return
	}

	voteToCreate.Thread = existingThread.ID
	voteToCreate.Nickname = existingUser.Nickname

	fmt.Println(voteToCreate)

	updatedThread, err, status := models.CreateVoteAndUpdateThread(voteToCreate)
	if err != nil {
		ErrResponse(res, status, err.Error())

		return
	}

	fmt.Println("--after vote--")
	PrintThread(updatedThread)

	ResponseObject(res, http.StatusOK, updatedThread)
}
