package models

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"time"
	"tp_db_forum/internal/database"
)

type Post struct {
	Author   string    `json:"author"`
	Created  time.Time `json:"created"`
	Forum    string    `json:"forum"`
	ID       int64     `json:"id"`
	IsEdited bool      `json:"isEdited"`
	Message  string    `json:"message"`
	Parent   int64     `json:"parent"`
	Thread   int32     `json:"thread"`
}

func CreatePosts(postsToCreate []Post, existingThread Thread) ([]Post, error, int) {
	conn := database.Connection

	now := time.Now()

	for i, post := range postsToCreate {
		fmt.Println("--", i, "--")

		resInsert, err := conn.Exec(`INSERT INTO forum_post (author, created, forum, message, parent, thread) VALUES ($1, $2, $3, $4, $5, $6)`,
			post.Author, now, existingThread.Forum, post.Message, post.Parent, existingThread.ID)

		if resInsert.RowsAffected() == 0 {
			return []Post{}, errors.Wrap(err, "cant create thread"), http.StatusInternalServerError
		}
		postsToCreate[i].Forum = existingThread.Forum
		postsToCreate[i].Thread = existingThread.ID
		postsToCreate[i].Created = now

		// get last id
		res, err := conn.Query(`SELECT last_value FROM forum_post_id_seq`)
		for res.Next() {
			err := res.Scan(&postsToCreate[i].ID)

			if err != nil {
				return []Post{}, errors.Wrap(err, "db query result parsing error"), http.StatusInternalServerError
			}
		}
	}
	return postsToCreate, nil, http.StatusOK
}
