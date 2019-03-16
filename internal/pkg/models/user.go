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

		return u, nil
	}

	return User{}, errors.New("cannot get user by nickname")
}

func GetUserByEmail(email string) (User, error) {
	conn := database.Connection
	u := User{}
	res, err := conn.Query(`SELECT about, email, fullname, nickname FROM forum_users WHERE email = $1`, email)
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

func GetUserByNicknameOrEmail(nickname string, email string) ([]User, error) {
	result := make([]User, 0, 1)
	conn := database.Connection
	u := User{}
	res, err := conn.Query(`SELECT about, email, fullname, nickname FROM forum_users WHERE email = $1 OR nickname = $2`,
		email, nickname)
	if err != nil {
		return []User{}, errors.Wrap(err, "cannot get user by nickname or email")
	}

	if res.Next() {
		err := res.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname)
		if err != nil {
			return []User{}, errors.Wrap(err, "db query result parsing error")
		}
		result = append(result, u)
	}

	return result, nil
}

func CreateUser(userToCreate User) (User, error) {
	conn := database.Connection

	_, err := conn.Exec(`Insert Into forum_users (nickname, fullname, email, about) VALUES ($1, $2, $3, $4)`,
		userToCreate.Nickname, userToCreate.Fullname, userToCreate.Email, userToCreate.About)
	if err != nil {
		return User{}, errors.Wrap(err, "cannot create user")
	}

	return userToCreate, nil
}

func UpdateUser(userToUpdate User) (User, error) {
	conn := database.Connection

	_, err := conn.Exec(`Update forum_users SET fullname = $1, email = $2, about = $3 WHERE nickname = $4`,
		userToUpdate.Fullname, userToUpdate.Email, userToUpdate.About, userToUpdate.Nickname)
	if err != nil {
		return User{}, errors.Wrap(err, "cannot update user")
	}

	return userToUpdate, nil
}
