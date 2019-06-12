package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Smet1/tp_db_forum/internal/pkg/models"
	"github.com/pkg/errors"
)

func PrintThread(t models.Thread) {
	fmt.Println("\tauthor = ", t.Author)
	fmt.Println("\tcreated = ", t.Created)
	fmt.Println("\tforum = ", t.Forum)
	fmt.Println("\tid = ", t.ID)
	fmt.Println("\tmessage = ", t.Message)
	fmt.Println("\tslug = ", t.Slug)
	fmt.Println("\ttitle = ", t.Title)
	fmt.Println("\tvotes = ", t.Votes)
}

func CreateThread(res http.ResponseWriter, req *http.Request) {
	//log.Println("CreateThread", req.URL)

	slugName, _ := checkVar("slug", req)
	//if err != nil {
	//	ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get user slug").Error())
	//	return
	//}

	t := models.Thread{}
	body, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	_ = t.UnmarshalJSON(body)

	t.Forum = slugName.(string)

	createdThread, err, status := models.CreateThread(t)
	if err != nil {
		if status == http.StatusNotFound {
			ErrResponse(res, status, err.Error())

			return
		}

		if status == http.StatusConflict {
			conflictThread, _, _ := models.GetThreadByIDorSlug(-1, t.Slug)
			//ResponseObject(res, status, conflictThread)
			ResponseEasyObject(res, status, conflictThread)

			return
		}

		//existingForum, err := models.GetForumBySlug(f.Slug)
		//if err != nil {
		//	ErrResponse(res, status, err.Error())
		//	return
		//}
		//
		//ResponseObject(res, http.StatusConflict, existingForum)
		//return
	}

	ResponseEasyObject(res, http.StatusCreated, createdThread)
	//ResponseObject(res, http.StatusCreated, createdThread)
}

func GetThreads(res http.ResponseWriter, req *http.Request) {
	//log.Println("GetThreads", req.URL)

	slugName, _ := checkVar("slug", req)
	//if err != nil {
	//	ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get user slug").Error())
	//	return
	//}

	query := req.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	since := query.Get("since")
	desc, _ := strconv.ParseBool(query.Get("desc"))

	threads, err, status := models.GetForumThreads(slugName.(string), limit, since, desc)
	if err != nil {
		if status == http.StatusNotFound {
			ErrResponse(res, status, err.Error())

			return
		}

		ErrResponse(res, status, err.Error())

		return
	}

	//ResponseObject(res, http.StatusOK, threads)
	ResponseEasyObject(res, http.StatusOK, threads)
}

func UpdateThread(res http.ResponseWriter, req *http.Request) {
	//log.Println("UpdateThread", req.URL)

	slugOrId, _ := checkVar("slug_or_id", req)
	//if err != nil {
	//	ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get slug or id").Error())
	//	return
	//}

	slug := slugOrId.(string)
	id, _ := strconv.ParseInt(slug, 10, 32)
	if id == 0 {
		id = -1
	}

	t := models.Thread{}
	body, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	t.UnmarshalJSON(body)

	existingThread, err, status := models.GetThreadByIDorSlug(int(id), slug)
	if err != nil {
		ErrResponse(res, status, errors.Wrap(err, "not found").Error())

		return
	}

	updatedThread, err, status := models.UpdateThread(existingThread, t)
	if err != nil {
		ErrResponse(res, status, err.Error())

		return
	}

	//ResponseObject(res, status, updatedThread)
	ResponseEasyObject(res, status, updatedThread)
}
