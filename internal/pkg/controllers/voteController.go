package controllers

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"math/rand"
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

		log.Println("CreateVote", err)
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

		log.Println("CreateVote:", errors.Wrap(err, "slug not found").Error(), "id = ", id, "slug = ", slug)
		return
	}

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

	updatedThread, err, status := models.CreateVoteAndUpdateThread(voteToCreate)
	if err != nil {
		ErrResponse(res, status, err.Error())

		return
	}

	idLog := rand.Int31n(1000)
	fmt.Println("--VOTE-- idLog=", idLog)
	fmt.Println(voteToCreate)
	fmt.Println("--before vote-- idLog=", idLog)
	PrintThread(existingThread)
	fmt.Println("--after vote-- idLog=", idLog)
	PrintThread(updatedThread)

	ResponseObject(res, http.StatusOK, updatedThread)
}

func GetThreadDetails(res http.ResponseWriter, req *http.Request) {
	log.Println("=============")
	log.Println("GetThreadDetails", req.URL)

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

	ResponseObject(res, http.StatusOK, existingThread)
}
