package models

import (
	"github.com/pkg/errors"
	"tp_db_forum/internal/database"
)

type User struct {
	About    string `json:"about"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname"`
}

func GetUserByNickname(nickname string) (User, error) {
	conn := database.Connection
	u := User{}
	res, err := conn.Query(`SELECT about, email, fullname, nickname FROM forum_users WHERE nickname = $1`, nickname)
	if err != nil {
		return User{}, errors.Wrap(err, "cannot get user by nickname")
	}

	if res.Next() {
		err := res.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname)
		if err != nil {
			return User{}, errors.Wrap(err, "db query result parsing error")
		}
	}

	return u, nil
}
