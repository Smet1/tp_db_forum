package models

import (
	"github.com/pkg/errors"
	"net/http"
	"tp_db_forum/internal/database"
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
	res, err := conn.Query("SELECT * FROM (SELECT count(posts) FROM forum_forum) as ff" +
		" CROSS JOIN (SELECT count(id) FROM forum_post) as fp" +
		" CROSS JOIN (SELECT count(id) FROM forum_thread) as ft" +
		" CROSS JOIN (SELECT count(id) FROM forum_users) as fu")
	defer res.Close()

	if err != nil {
		return Status{}, errors.Wrap(err, "cant get db statisics"), http.StatusInternalServerError
	}

	s := Status{}
	for res.Next() {
		err := res.Scan(&s.Forum, &s.Post, &s.Thread, &s.User)

		if err != nil {
			return Status{}, errors.Wrap(err, "db query result parsing error"), http.StatusInternalServerError
		}

		return s, nil, http.StatusOK
	}

	return Status{}, errors.Wrap(err, "cant get db statisics"), http.StatusInternalServerError
}

func ClearDB() (error, int) {
	conn := database.Connection

	res, err := conn.Query("TRUNCATE TABLE forum_users, forum_forum, forum_thread, forum_post, forum_vote CASCADE")
	defer res.Close()

	if err != nil {
		return errors.Wrap(err, "cant truncate db"), http.StatusInternalServerError
	}

	return nil, http.StatusOK
}
