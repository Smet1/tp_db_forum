package models

import (
	"github.com/pkg/errors"
	"tp_db_forum/internal/database"
)

type Forum struct {
	Posts   int64  `json:"posts"`
	Slug    string `json:"slug"`
	Threads int32  `json:"threads"`
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

	_, err = conn.Exec(`Insert Into forum_forum (posts, slug, threads, title, "user") VALUES ($1, $2, $3, $4, $5)`,
		forumToCreate.Posts, forumToCreate.Slug, forumToCreate.Threads, forumToCreate.Title, forumToCreate.User)
	if err != nil {
		return Forum{}, errors.Wrap(err, "cannot create forum")
	}

	return forumToCreate, nil
}

func GetForumBySlug(slug string) (Forum, error) {
	conn := database.Connection
	res, err := conn.Query(`SELECT * FROM forum_forum WHERE slug = $1`, slug)
	if err != nil {
		return Forum{}, errors.Wrap(err, "cannot get forum by slug")
	}

	f := Forum{}

	if res.Next() {
		err := res.Scan(&f.Posts, &f.Slug, &f.Threads, &f.Title, &f.User)
		if err != nil {
			return Forum{}, errors.Wrap(err, "db query result parsing error")
		}
	}

	return f, nil
}
