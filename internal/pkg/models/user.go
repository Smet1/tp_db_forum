package models

import (
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"tp_db_forum/internal/database"
)

//easyjson:json
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
	defer res.Close()

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
	defer res.Close()

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
	defer res.Close()

	if err != nil {
		return []User{}, errors.Wrap(err, "cannot get user by nickname or email")
	}

	for res.Next() {
		err := res.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname)

		if err != nil {
			return []User{}, errors.Wrap(err, "db query result parsing error")
		}
		result = append(result, u)
	}

	//fmt.Println("==========")
	//for _, val := range result {
	//	fmt.Println("\t", val)
	//}
	//fmt.Println("==========")

	return result, nil
}

func CreateUser(userToCreate User) (User, error) {
	conn := database.Connection

	res, err := conn.Exec(`Insert Into forum_users (nickname, fullname, email, about) VALUES ($1, $2, $3, $4)`,
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

	baseSql := "Update forum_users SET"
	if userToUpdate.Fullname == "" {
		baseSql += " fullname = fullname,"
	} else {
		baseSql += " fullname = '" + userToUpdate.Fullname + "',"
	}

	if userToUpdate.Email == "" {
		baseSql += " email = email,"
	} else {
		baseSql += " email = '" + userToUpdate.Email + "',"
	}

	if userToUpdate.About == "" {
		baseSql += " about = about"
	} else {
		baseSql += " about = '" + userToUpdate.About + "'"
	}

	baseSql += " WHERE nickname = '" + userToUpdate.Nickname + "'"

	//fmt.Println(baseSql)

	res, err := conn.Exec(baseSql)
	if err != nil {
		return User{}, errors.Wrap(err, "cannot update user"), http.StatusConflict
	}

	if res.RowsAffected() == 0 {
		return User{}, errors.New("not found"), http.StatusNotFound
	}

	updatedUser, _ := GetUserByNickname(userToUpdate.Nickname)

	return updatedUser, nil, http.StatusOK
}

func GetForumUsersBySlug(existingForum Forum, limit int, since string, desc bool) ([]User, error, int) {
	conn := database.Connection

	//select DISTINCT about, email, fullname, nickname
	//FROM forum_users
	//       LEFT JOIN forum_post fp ON fp.author = nickname AND nickname > 'N1HR3DVSHJZhr.bill'
	//       LEFT JOIN forum_thread ft ON ft.author = nickname AND nickname > 'N1HR3DVSHJZhr.bill'
	//WHERE fp.forum = 'xE6RM2vYIkOoK'
	//   OR ft.forum = 'xE6RM2vYIkOoK'
	//ORDER BY nickname
	//LIMIT 4;

	baseSQL := `SELECT DISTINCT about, email, fullname, nickname FROM forum_users`
	if since != "" {
		if desc {
			baseSQL += ` LEFT JOIN forum_post fp ON fp.author = nickname` + " AND nickname < '" + since + "'"
			baseSQL += ` LEFT JOIN forum_thread ft ON ft.author = nickname` + " AND nickname < '" + since + "'"
		} else {
			baseSQL += ` LEFT JOIN forum_post fp ON fp.author = nickname` + " AND nickname > '" + since + "'"
			baseSQL += ` LEFT JOIN forum_thread ft ON ft.author = nickname` + " AND nickname > '" + since + "'"

		}
	} else {
		baseSQL += ` LEFT JOIN forum_post fp ON fp.author = nickname`
		baseSQL += ` LEFT JOIN forum_thread ft ON ft.author = nickname`
	}
	baseSQL += ` WHERE fp.forum = '` + existingForum.Slug + `' OR ft.forum = '` + existingForum.Slug + `'`

	if desc {
		baseSQL += " ORDER BY nickname DESC"
	} else {
		baseSQL += " ORDER BY nickname ASC"
	}

	if limit != 0 {
		baseSQL += " LIMIT " + strconv.Itoa(limit)
	}

	//fmt.Println("\t", baseSQL)

	res, err := conn.Query(baseSQL)
	if err != nil {
		return []User{}, errors.Wrap(err, "cannot get user by nickname or email"), http.StatusInternalServerError
	}
	defer res.Close()

	queriedUsers := make([]User, 0, 1)
	u := User{}

	for res.Next() {
		err := res.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname)

		if err != nil {
			return []User{}, errors.Wrap(err, "db query result parsing error"), http.StatusInternalServerError
		}
		queriedUsers = append(queriedUsers, u)
	}

	return queriedUsers, nil, http.StatusOK

}
