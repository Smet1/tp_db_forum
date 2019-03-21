package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

func checkVar(varName string, req *http.Request) (interface{}, error) {
	requestVariables := mux.Vars(req)
	if requestVariables == nil {

		return nil, errors.New("user nickname not provided")
	}

	result, ok := requestVariables[varName]
	if !ok {

		return nil, errors.New("vars found, but cant found nickname")
	}

	return result, nil
}

func GetUserProfile(res http.ResponseWriter, req *http.Request) {
	log.Println("=============")
	log.Println("GetUserProfile", req.URL)

	searchingNickname, err := checkVar("nickname", req)
	if err != nil {
		ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get user nickname").Error())
		return
	}

	u, err := models.GetUserByNickname(searchingNickname.(string))
	if err != nil || u.Email == "" {
		ErrResponse(res, http.StatusNotFound, "Can't find user")
		return
	}

	ResponseObject(res, http.StatusOK, u)
}

func UpdateUserProfile(res http.ResponseWriter, req *http.Request) {
	log.Println("=============")
	log.Println("UpdateUserProfile", req.URL)

	nicknameToUpdate, err := checkVar("nickname", req)
	if err != nil {
		ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get user nickname").Error())
		return
	}

	u := models.User{}

	//body, _ := ioutil.ReadAll(req.Body)
	//u.UnmarshalJSON(body)

	status, err := ParseRequestIntoStruct(req, &u)
	if err != nil {
		ErrResponse(res, status, err.Error())

		log.Println("\t", errors.Wrap(err, "ParseRequestIntoStruct error"))
		return
	}
	u.Nickname = nicknameToUpdate.(string)

	updatedUser, err, errCode := models.UpdateUser(u)
	if err != nil {
		ErrResponse(res, errCode, err.Error())
		return
	}

	ResponseObject(res, http.StatusOK, updatedUser)
}

func CreateUser(res http.ResponseWriter, req *http.Request) {
	log.Println("=============")
	log.Println("CreateUser", req.URL)

	nicknameToCreate, err := checkVar("nickname", req)
	if err != nil {
		ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get user nickname").Error())
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

	u.Nickname = nicknameToCreate.(string)
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

	ResponseObject(res, http.StatusCreated, createdUser)
}

func GetForumUsers(res http.ResponseWriter, req *http.Request) {
	log.Println("=============")
	log.Println("GetForumUsers", req.URL)

	query := req.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	since := query.Get("since")
	desc, _ := strconv.ParseBool(query.Get("desc"))

	fmt.Println(query)
	fmt.Println(limit)
	fmt.Println(since)
	fmt.Println(desc)

	searchingSlug, err := checkVar("slug", req)
	if err != nil {
		ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get forum slug").Error())
		return
	}

	fmt.Println(searchingSlug)

	existingForum, err := models.GetForumBySlug(searchingSlug.(string))
	if err != nil {
		ErrResponse(res, http.StatusNotFound, errors.Wrap(err, "not found").Error())
		return
	}

	users, err, status := models.GetForumUsersBySlug(existingForum, limit, since, desc)
	if err != nil {
		if status == http.StatusNotFound {
			ErrResponse(res, status, err.Error())
			return
		}

		ErrResponse(res, status, err.Error())
		return
	}

	ResponseObject(res, http.StatusOK, users)
}
