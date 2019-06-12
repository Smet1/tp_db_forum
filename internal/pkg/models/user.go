package models

import (
	"net/http"
	"strconv"

	"github.com/Smet1/tp_db_forum/internal/database"
	"github.com/pkg/errors"
)

//easyjson:json
type User struct {
	About    string `json:"about,omitempty"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname,omitempty"`
}

//easyjson:json
type Users []User

func GetUserByNickname(nickname string) (User, error) {
	conn := database.Connection
	u := User{}
	res, err := conn.Query(`SELECT about, email, fullname, nickname FROM forum_users WHERE nickname = $1`, nickname)
	if err != nil {
		return User{}, errors.Wrap(err, "cannot get user by nickname")
	}
	defer res.Close()

	if res.Next() {
		_ = res.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname)
		//if err != nil {
		//	return User{}, errors.Wrap(err, "db query result parsing error")
		//}

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
	defer res.Close()

	if res.Next() {
		err := res.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname)
		if err != nil {
			return User{}, errors.Wrap(err, "db query result parsing error")
		}
	}

	return u, nil
}

func GetUserByNicknameOrEmail(nickname string, email string) (Users, error) {
	result := make([]User, 0, 1)
	conn := database.Connection
	u := User{}
	res, err := conn.Query(`SELECT about, email, fullname, nickname FROM forum_users WHERE email = $1 OR nickname = $2`,
		email, nickname)
	if err != nil {
		return []User{}, errors.Wrap(err, "cannot get user by nickname or email")
	}
	defer res.Close()

	for res.Next() {
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

	res, err := conn.Exec(`INSERT INTO forum_users (nickname, fullname, email, about) VALUES ($1, $2, $3, $4)`,
		userToCreate.Nickname, userToCreate.Fullname, userToCreate.Email, userToCreate.About)
	if err != nil {
		return User{}, errors.Wrap(err, "cannot create user")
	}

	if res.RowsAffected() == 0 {
		return User{}, errors.Wrap(err, "cannot create user")
	}

	return userToCreate, nil
}

func UpdateUser(userToUpdate User) (User, error, int) {
	conn := database.Connection

	if userToUpdate.About == "" && userToUpdate.Email == "" && userToUpdate.Fullname == "" {
		updatedUser, _ := GetUserByNickname(userToUpdate.Nickname)

		return updatedUser, nil, http.StatusOK
	}

	baseSQL := "Update forum_users SET"
	if userToUpdate.Fullname == "" {
		baseSQL += " fullname = fullname,"
	} else {
		baseSQL += " fullname = '" + userToUpdate.Fullname + "',"
	}

	if userToUpdate.Email == "" {
		baseSQL += " email = email,"
	} else {
		baseSQL += " email = '" + userToUpdate.Email + "',"
	}

	if userToUpdate.About == "" {
		baseSQL += " about = about"
	} else {
		baseSQL += " about = '" + userToUpdate.About + "'"
	}

	baseSQL += " WHERE nickname = '" + userToUpdate.Nickname + "'"

	res, err := conn.Exec(baseSQL)
	if err != nil {
		return User{}, errors.Wrap(err, "cannot update user"), http.StatusConflict
	}

	if res.RowsAffected() == 0 {
		return User{}, errors.New("not found"), http.StatusNotFound
	}

	updatedUser, _ := GetUserByNickname(userToUpdate.Nickname)

	return updatedUser, nil, http.StatusOK
}

func GetForumUsersBySlug(existingForum Forum, limit int, since string, desc bool) (Users, error, int) {
	conn := database.Connection

	baseSQL := `SELECT about, email, fullname, fu.nickname FROM forum_users_forum JOIN forum_users fu ON fu.nickname = forum_users_forum.nickname`

	baseSQL += ` where slug = '` + existingForum.Slug + `'`
	if since != "" {
		if desc {
			baseSQL += ` AND fu.nickname < '` + since + `'`
		} else {
			baseSQL += ` AND fu.nickname > '` + since + `'`
		}
	}

	if desc {
		baseSQL += " ORDER BY nickname DESC"
	} else {
		baseSQL += " ORDER BY nickname ASC"
	}

	if limit != 0 {
		baseSQL += " LIMIT " + strconv.Itoa(limit)
	}

	res, _ := conn.Query(baseSQL)
	//if err != nil {
	//	return []User{}, errors.Wrap(err, "cannot get user by nickname or email"), http.StatusInternalServerError
	//}
	defer res.Close()

	queriedUsers := make([]User, 0, 1)
	u := User{}

	for res.Next() {
		_ = res.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname)

		//if err != nil {
		//	return []User{}, errors.Wrap(err, "db query result parsing error"), http.StatusInternalServerError
		//}
		queriedUsers = append(queriedUsers, u)
	}

	return queriedUsers, nil, http.StatusOK

}

func AddUser(nickname string, forumSlug string) {
	conn := database.Connection

	_, _ = conn.Exec(`INSERT INTO forum_users_forum (nickname, slug) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		nickname, forumSlug)

	return
}
