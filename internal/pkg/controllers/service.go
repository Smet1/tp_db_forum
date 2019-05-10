package controllers

import (
	"github.com/Smet1/tp_db_forum/internal/pkg/models"
	"net/http"
)

func GetDBStatus(res http.ResponseWriter, req *http.Request) {
	//log.Println("=============")
	//log.Println("GetDBStatus", req.URL)

	stats, err, status := models.GetDBCountData()
	if err != nil {
		ErrResponse(res, status, err.Error())

		return
	}

	ResponseObject(res, status, stats)
}

func ClearDB(res http.ResponseWriter, req *http.Request) {
	//log.Println("=============")
	//log.Println("ClearDB", req.URL)

	err, status := models.ClearDB()
	if err != nil {
		ErrResponse(res, status, err.Error())
		return
	}

	OkResponse(res, "clear ok")
}
