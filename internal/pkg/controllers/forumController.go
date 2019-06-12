package controllers

import (
	"io/ioutil"
	"net/http"

	"github.com/Smet1/tp_db_forum/internal/pkg/models"
)

func CreateForum(res http.ResponseWriter, req *http.Request) {
	//log.Println("CreateForum", req.URL)

	f := models.Forum{}
	body, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	_ = f.UnmarshalJSON(body)

	createdForum, err := models.CreateForum(f)
	if err != nil {
		_, err := models.GetUserByNickname(f.User)
		if err != nil {
			ErrResponse(res, http.StatusNotFound, "Can't find user with nickname "+f.User)
			return
		}

		existingForum, err := models.GetForumBySlug(f.Slug)
		if err != nil {
			ErrResponse(res, http.StatusNotFound, err.Error())
			return
		}

		//ResponseObject(res, http.StatusConflict, existingForum)
		ResponseEasyObject(res, http.StatusConflict, existingForum)
		return
	}

	//ResponseObject(res, http.StatusCreated, createdForum)
	ResponseEasyObject(res, http.StatusCreated, createdForum)
}

func GetForum(res http.ResponseWriter, req *http.Request) {
	//log.Println("GetForum", req.URL)

	searchingSlug, _ := checkVar("slug", req)
	//if err != nil {
	//	ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get forum slug").Error())
	//	return
	//}

	f, err := models.GetForumBySlug(searchingSlug.(string))
	if err != nil || f.User == "" {
		ErrResponse(res, http.StatusNotFound, "Can't find slug")
		return
	}

	//ResponseObject(res, http.StatusOK, f)
	ResponseEasyObject(res, http.StatusOK, f)
}
