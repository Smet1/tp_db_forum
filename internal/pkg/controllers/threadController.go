package controllers

import (
	"net/http"
)

func CreateThread(res http.ResponseWriter, req *http.Request) {
	//nicknameToCreate, err := checkVar("nickname", req)
	//if err != nil {
	//	ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get user nickname").Error())
	//	return
	//}
	//
	////u, err := models.GetUserByNickname(searchingNickname)
	////if err == nil && u.Email != "" {
	////	ErrResponseObject(res, http.StatusConflict, u)
	////
	////	return
	////}
	//
	//u := models.User{}
	//status, err := ParseRequestIntoStruct(req, &u)
	//if err != nil {
	//	ErrResponse(res, status, err.Error())
	//
	//	log.Println("\t", errors.Wrap(err, "ParseRequestIntoStruct error"))
	//	return
	//}
	//
	//u.Nickname = nicknameToCreate.(string)
	////existingUser, err := models.GetUserByEmail(u.Email)
	////if err == nil && u.Email != "" {
	////	ErrResponseObject(res, http.StatusConflict, existingUser)
	////
	////	return
	////}
	//
	//createdUser, err := models.CreateUser(u)
	//if err != nil {
	//	exitingUsers, err := models.GetUserByNicknameOrEmail(u.Nickname, u.Email)
	//	if err != nil {
	//		ErrResponse(res, status, err.Error())
	//		return
	//	}
	//
	//	ResponseObject(res, http.StatusConflict, exitingUsers)
	//	return
	//}
	//
	//ResponseObject(res, http.StatusCreated, createdUser)
}
