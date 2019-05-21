package controllers

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/Smet1/tp_db_forum/internal/pkg/models"
	"github.com/pkg/errors"
)

func CreatePosts(res http.ResponseWriter, req *http.Request) {
	//log.Println("CreatePosts", req.URL)

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
		ErrResponse(res, status, errors.Wrap(err, "not found").Error())
		return
	}

	//postsToCreate := make([]models.Post, 0, 1)
	//status, err = ParseRequestIntoStruct(req, &postsToCreate)
	//if err != nil {
	//	ErrResponse(res, status, err.Error())
	//
	//	//log.Println("\t", errors.Wrap(err, "ParseRequestIntoStruct error"))
	//	return
	//}

	postsToCreate := models.Posts{}
	body, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	_ = postsToCreate.UnmarshalJSON(body)
	if len(postsToCreate) == 0 {
		ResponseObject(res, http.StatusCreated, postsToCreate)
		return
	}

	createdPosts, err, status := models.CreatePosts(postsToCreate, existingThread)
	if err != nil {
		ErrResponse(res, status, err.Error())

		//log.Println("\t", errors.Wrap(err, "models.CreatePost error"))
		return
	}

	ResponseObject(res, http.StatusCreated, createdPosts)
}

func GetThreadPosts(res http.ResponseWriter, req *http.Request) {
	//log.Println("GetThreadPosts", req.URL)

	slugOrID, _ := checkVar("slug_or_id", req)
	//if err != nil {
	//	ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get user slug").Error())
	//	return
	//}
	slug := slugOrID.(string)
	id, err := strconv.ParseInt(slug, 10, 32)
	if id == 0 {
		id = -1
	}

	existingThread, err, status := models.GetThreadByIDorSlug(int(id), slug)
	if err != nil {
		ErrResponse(res, status, errors.Wrap(err, "not found").Error())
		return
	}

	query := req.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	since, _ := strconv.Atoi(query.Get("since"))
	sort := query.Get("sort")
	desc, _ := strconv.ParseBool(query.Get("desc"))

	sortedPosts, err, status := models.GetSortedPosts(existingThread, limit, since, sort, desc)
	if err != nil {
		ErrResponse(res, status, err.Error())

		return
	}

	ResponseObject(res, status, sortedPosts)
}

func UpdatePost(res http.ResponseWriter, req *http.Request) {
	//log.Println("UpdatePost", req.URL)

	postId, _ := checkVar("id", req)
	//if err != nil {
	//	ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get post id").Error())
	//	return
	//}
	id, err := strconv.ParseInt(postId.(string), 10, 64)
	if id == 0 {
		id = -1
	}

	existingPost, err, status := models.GetPostByID(id)
	if err != nil {
		ErrResponse(res, status, err.Error())

		return
	}

	newPost := models.Post{}
	body, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	_ = newPost.UnmarshalJSON(body)

	updatedPost, err, status := models.UpdatePost(existingPost, newPost)

	if err != nil {
		ErrResponse(res, status, err.Error())

		return
	}

	ResponseObject(res, status, updatedPost)
}

func GetPostInfo(res http.ResponseWriter, req *http.Request) {
	//log.Println("GetPostInfo", req.URL)

	slug, _ := checkVar("id", req)
	//if err != nil {
	//	ErrResponse(res, http.StatusBadRequest, errors.Wrap(err, "cant get post id").Error())
	//	return
	//}
	id, _ := strconv.ParseInt(slug.(string), 10, 64)
	if id == 0 {
		id = -1
	}

	query := req.URL.Query()

	related := strings.Split(query.Get("related"), ",")

	existingPost, err, status := models.GetPostByID(id)
	if err != nil {
		ErrResponse(res, status, err.Error())

		return
	}

	tupaKek, err, status := models.GetPostDetails(existingPost, related)
	if err != nil {
		ErrResponse(res, status, err.Error())

		return
	}

	tupaKek.Post = &existingPost

	ResponseObject(res, status, tupaKek)
}
