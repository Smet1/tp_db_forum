package models

import (
	"github.com/Smet1/tp_db_forum/internal/database"
	"github.com/pkg/errors"
	"net/http"
)

//easyjson:json
type Status struct {
	Forum  int32 `json:"forum"`
	Post   int32 `json:"post"`
	Thread int32 `json:"thread"`
	User   int32 `json:"user"`
}

func GetDBCountData() (Status, error, int) {
	conn := database.Connection
	res, _ := conn.Query("SELECT * FROM (SELECT count(posts) FROM forum_forum) as ff" +
		" CROSS JOIN (SELECT count(id) FROM forum_post) as fp" +
		" CROSS JOIN (SELECT count(id) FROM forum_thread) as ft" +
		" CROSS JOIN (SELECT count(nickname) FROM forum_users) as fu")
	defer res.Close()

	//if err != nil {
	//	return Status{}, errors.Wrap(err, "cant get db statisics"), http.StatusInternalServerError
	//}

	s := Status{}
	for res.Next() {
		_ = res.Scan(&s.Forum, &s.Post, &s.Thread, &s.User)

		//if err != nil {
		//	return Status{}, errors.Wrap(err, "db query result parsing error"), http.StatusInternalServerError
		//}

		return s, nil, http.StatusOK
	}

	return Status{}, errors.New("cant get db statisics"), http.StatusInternalServerError
}

func ClearDB() (error, int) {
	conn := database.Connection

	res, _ := conn.Query("TRUNCATE TABLE forum_users, forum_forum, forum_thread, forum_post, forum_vote CASCADE")
	//if err != nil {
	//	return errors.Wrap(err, "cant truncate db"), http.StatusInternalServerError
	//}
	defer res.Close()

	return nil, http.StatusOK
}
