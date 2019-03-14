package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"tp_db_forum/internal/pkg/models"
)

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
	fmt.Println("POST --- create")
}
