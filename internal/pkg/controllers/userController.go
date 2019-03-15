package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"tp_db_forum/internal/pkg/models"
)

func ParseRequestIntoStruct(req *http.Request, requestStruct interface{}) (int, error) {

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "body parsing error")
	}

	err = json.Unmarshal(body, &requestStruct)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "json parsing error")
	}

	return 0, nil
}

func GetUserProfile(res http.ResponseWriter, req *http.Request) {
	requestVariables := mux.Vars(req)
	if requestVariables == nil {
		ErrResponse(res, http.StatusBadRequest, "user id not provided")

		log.Println("\t", errors.New("no vars found"))
		return
	}

	searchingNickname, ok := requestVariables["nickname"]
	if !ok {
		ErrResponse(res, http.StatusInternalServerError, "error")

		log.Println("\t", errors.New("vars found, but cant found nickname"))
		return
	}

	u, err := models.GetUserByNickname(searchingNickname)
	if err != nil || u.Email == "" {
		ErrResponse(res, http.StatusNotFound, "Can't find user")
		return
	}

	OkResponse(res, u)
}

func UpdateUserProfile(res http.ResponseWriter, req *http.Request) {

}

func CreateUser(res http.ResponseWriter, req *http.Request) {
	requestVariables := mux.Vars(req)
	if requestVariables == nil {
		ErrResponse(res, http.StatusBadRequest, "user nickname not provided")

		log.Println("\t", errors.New("no vars found"))
		return
	}

	nicknameToCreate, ok := requestVariables["nickname"]
	if !ok {
		ErrResponse(res, http.StatusInternalServerError, "error")

		log.Println("\t", errors.New("vars found, but cant found nickname"))
		return
	}

	//u, err := models.GetUserByNickname(searchingNickname)
	//if err == nil && u.Email != "" {
	//	ErrResponseObject(res, http.StatusConflict, u)
	//
	//	return
	//}

	u := models.User{}
	status, err := ParseRequestIntoStruct(req, &u)
	if err != nil {
		ErrResponse(res, status, err.Error())

		log.Println("\t", errors.Wrap(err, "ParseRequestIntoStruct error"))
		return
	}


	u.Nickname = nicknameToCreate
	//existingUser, err := models.GetUserByEmail(u.Email)
	//if err == nil && u.Email != "" {
	//	ErrResponseObject(res, http.StatusConflict, existingUser)
	//
	//	return
	//}

	createdUser, err := models.CreateUser(u)
	if err != nil {
		exitingUsers, err := models.GetUserByNicknameOrEmail(u.Nickname, u.Email)
		if err != nil {
			ErrResponse(res, status, err.Error())
			return
		}

		ResponseObject(res, http.StatusConflict, exitingUsers)
		return
	}

	ResponseObject(res, http.StatusCreated ,createdUser)
}
