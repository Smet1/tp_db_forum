package models

import (
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"tp_db_forum/internal/database"
)

//easyjson:json
type Forum struct {
	Posts   int64  `json:"posts,omitempty"`
	Slug    string `json:"slug"`
	Threads int32  `json:"threads,omitempty"`
	Title   string `json:"title"`
	User    string `json:"user"`
}

func CreateForum(forumToCreate Forum) (Forum, error) {
	conn := database.Connection

	existingUser, err := GetUserByNickname(forumToCreate.User)
	if err != nil {
		return Forum{}, errors.Wrap(err, "cannot create forum")
	}

	forumToCreate.User = existingUser.Nickname

	_, err = conn.Exec(`INSERT INTO forum_forum (posts, slug, threads, title, "user") VALUES ($1, $2, $3, $4, $5)`,
		forumToCreate.Posts, forumToCreate.Slug, forumToCreate.Threads, forumToCreate.Title, forumToCreate.User)
	if err != nil {
		return Forum{}, errors.Wrap(err, "cannot create forum")
	}

	//AddUser(forumToCreate.User, forumToCreate.Slug)

	return forumToCreate, nil
}

func GetForumBySlug(slug string) (Forum, error) {
	conn := database.Connection
	res, err := conn.Query(`SELECT * FROM forum_forum WHERE slug = $1`, slug)
	if err != nil {
		return Forum{}, errors.Wrap(err, "cannot get forum by slug")
	}
	defer res.Close()

	f := Forum{}

	if res.Next() {
		err := res.Scan(&f.Posts, &f.Slug, &f.Threads, &f.Title, &f.User)
		if err != nil {
			return Forum{}, errors.Wrap(err, "db query result parsing error")
		}
	}

	if f.Slug == "" {
		return Forum{}, errors.New("cannot get forum by slug")
	}

	return f, nil
}

func UpdateForumStats(existingForum Forum, varToUpdate string, inc bool, diff int) int {
	conn := database.Connection

	baseSQL := "UPDATE forum_forum SET"

	switch varToUpdate {
	case "post":
		baseSQL += " posts = posts"
		if inc {
			baseSQL += "+"
		} else {
			baseSQL += "-"
		}

		baseSQL += strconv.Itoa(diff)

	case "thread":
		baseSQL += " threads = threads"
		if inc {
			baseSQL += "+"
		} else {
			baseSQL += "-"
		}

		baseSQL += strconv.Itoa(diff)

	default:
		return http.StatusInternalServerError
	}

	baseSQL += " WHERE slug = $1"

	res, err := conn.Exec(baseSQL, existingForum.Slug)
	if err != nil {
		return http.StatusConflict
	}

	if res.RowsAffected() == 0 {
		return http.StatusNotFound
	}

	return http.StatusOK
}
