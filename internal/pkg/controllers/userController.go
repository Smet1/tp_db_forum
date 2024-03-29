package controllers

import (
	"goji.io/pat"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Smet1/tp_db_forum/internal/pkg/models"
	"github.com/pkg/errors"
)

func checkVar(varName string, req *http.Request) (string, error) {
	//requestVariables := mux.Vars(req)
	requestVariables := pat.Param(req, varName)
	//if requestVariables == nil {
	//
	//	return nil, errors.New("user nickname not provided")
	//}
	//
	//result, ok := requestVariables[varName]
	//if !ok {
	//
	//	return nil, errors.New("vars found, but cant found nickname")
	//}

	return requestVariables, nil
}

func GetUserProfile(res http.ResponseWriter, req *http.Request) {
	//log.Println("GetUserProfile", req.URL)

	searchingNickname, _ := checkVar("nickname", req)
	//if err != nil {
	//	ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get user nickname").Error())
	//	return
	//}

	u, err := models.GetUserByNickname(searchingNickname)
	if err != nil || u.Email == "" {
		ErrResponse(res, http.StatusNotFound, "Can't find user")

		return
	}

	//ResponseObject(res, http.StatusOK, u)
	ResponseEasyObject(res, http.StatusOK, u)
}

func UpdateUserProfile(res http.ResponseWriter, req *http.Request) {
	//log.Println("UpdateUserProfile", req.URL)

	nicknameToUpdate, _ := checkVar("nickname", req)
	//if err != nil {
	//	ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get user nickname").Error())
	//	return
	//}

	u := models.User{}
	body, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	u.UnmarshalJSON(body)

	u.Nickname = nicknameToUpdate

	updatedUser, err, errCode := models.UpdateUser(u)
	if err != nil {
		ErrResponse(res, errCode, err.Error())

		return
	}

	//ResponseObject(res, http.StatusOK, updatedUser)
	ResponseEasyObject(res, http.StatusOK, updatedUser)
}

func CreateUser(res http.ResponseWriter, req *http.Request) {
	//log.Println("CreateUser", req.URL)

	nicknameToCreate, _ := checkVar("nickname", req)
	//if err != nil {
	//	ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get user nickname").Error())
	//	return
	//}

	u := models.User{}
	body, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	_ = u.UnmarshalJSON(body)

	u.Nickname = nicknameToCreate

	createdUser, err := models.CreateUser(u)
	if err != nil {
		exitingUsers, err := models.GetUserByNicknameOrEmail(u.Nickname, u.Email)
		if err != nil {
			ErrResponse(res, http.StatusConflict, err.Error())

			return
		}

		//ResponseObject(res, http.StatusConflict, exitingUsers)
		ResponseEasyObject(res, http.StatusConflict, exitingUsers)

		return
	}

	ResponseEasyObject(res, http.StatusCreated, createdUser)
	//ResponseObject(res, http.StatusCreated, createdUser)
}

func GetForumUsers(res http.ResponseWriter, req *http.Request) {
	//log.Println("GetForumUsers", req.URL)

	query := req.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	since := query.Get("since")
	desc, _ := strconv.ParseBool(query.Get("desc"))

	searchingSlug, _ := checkVar("slug", req)
	//if err != nil {
	//	ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get forum slug").Error())
	//	return
	//}

	existingForum, err := models.GetForumBySlug(searchingSlug)
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

	//ResponseObject(res, http.StatusOK, users)
	ResponseEasyObject(res, http.StatusOK, users)
}
