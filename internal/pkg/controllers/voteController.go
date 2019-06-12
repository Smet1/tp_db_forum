package controllers

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Smet1/tp_db_forum/internal/pkg/models"
	"github.com/pkg/errors"
)

func CreateVote(res http.ResponseWriter, req *http.Request) {
	//log.Println("CreateVote", req.URL)

	slugOrID, _ := checkVar("slug_or_id", req)
	//if err != nil {
	//	ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get user slug").Error())
	//
	//	//log.Println("CreateVote", err)
	//	return
	//}
	slug := slugOrID.(string)
	id, _ := strconv.ParseInt(slug, 10, 32)
	if id == 0 {
		id = -1
	}

	existingThread, err, status := models.GetThreadByIDorSlug(int(id), slug)
	if err != nil {
		ErrResponse(res, status, errors.Wrap(err, "slug not found").Error())

		return
	}

	voteToCreate := models.Vote{}
	body, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	voteToCreate.UnmarshalJSON(body)

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

	//ResponseObject(res, http.StatusOK, updatedThread)
	ResponseEasyObject(res, http.StatusOK, updatedThread)
}

func GetThreadDetails(res http.ResponseWriter, req *http.Request) {
	//log.Println("GetThreadDetails", req.URL)

	slugOrID, _ := checkVar("slug_or_id", req)
	//if err != nil {
	//	ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get user slug").Error())
	//	return
	//}
	slug := slugOrID.(string)
	id, _ := strconv.ParseInt(slug, 10, 32)
	if id == 0 {
		id = -1
	}

	existingThread, err, status := models.GetThreadByIDorSlug(int(id), slug)
	if err != nil {
		ErrResponse(res, status, errors.Wrap(err, "slug not found").Error())

		return
	}

	//ResponseObject(res, http.StatusOK, existingThread)
	ResponseEasyObject(res, http.StatusOK, existingThread)
}
