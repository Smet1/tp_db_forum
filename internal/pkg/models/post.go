package models

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
	"tp_db_forum/internal/database"
)

type Post struct {
	Author   string    `json:"author"`
	Created  time.Time `json:"created,omitempty"`
	Forum    string    `json:"forum,omitempty"`
	ID       int64     `json:"id,omitempty"`
	IsEdited bool      `json:"isEdited,omitempty"`
	Message  string    `json:"message"`
	Parent   int64     `json:"parent,omitempty"`
	Thread   int32     `json:"thread,omitempty"`
	Path     []int32   `json:"-"`
}

func PrintPost(t Post) {
	fmt.Println("\tauthor = ", t.Author)
	fmt.Println("\tcreated = ", t.Created)
	fmt.Println("\tforum = ", t.Forum)
	fmt.Println("\tid = ", t.ID)
	fmt.Println("\tisEdited = ", t.IsEdited)
	fmt.Println("\tmessage = ", t.Message)
	fmt.Println("\tparent = ", t.Parent)
	fmt.Println("\tthread = ", t.Thread)
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

		PrintPost(postsToCreate[i])
	}
	return postsToCreate, nil, http.StatusOK
}

func GetSortedPosts(parentThread Thread, limit int, since int, sort string, desc bool) ([]Post, error, int) {
	// tree && parent tree UNDONE
	conn := database.Connection

	baseSQL := ""
	sortedPosts := make([]Post, 0, 1)
	strID := strconv.FormatInt(int64(parentThread.ID), 10)

	switch sort {
	case "flat":
		// author, created, forum, isedited, message, parent, thread
		baseSQL = "SELECT author, created, forum, id, isedited, message, parent, thread FROM forum_post WHERE thread = " + strID

		if since != 0 {
			if desc {
				baseSQL += " AND id < " + strconv.Itoa(since)
			} else {
				baseSQL += " AND id > " + strconv.Itoa(since)
			}
		}

		if desc {
			baseSQL += " ORDER BY id DESC"
		} else {
			baseSQL += " ORDER BY id"
		}

		baseSQL += " LIMIT " + strconv.Itoa(limit)
	case "tree":
		baseSQL = "SELECT author, created, forum, id, isedited, message, parent, thread FROM forum_post WHERE thread = " + strID

		if since != 0 {
			if desc {
				baseSQL += " AND id < (SELECT path FROM post WHERE id = " + strconv.Itoa(since) + ")"
			} else {
				baseSQL += " AND id > (SELECT path FROM post WHERE id = " + strconv.Itoa(since) + ")"
			}
		}

		if desc {
			baseSQL += " ORDER BY path DESC"
		} else {
			baseSQL += " ORDER BY path"
		}

		baseSQL += " LIMIT " + strconv.Itoa(limit)
	}

	res, err := conn.Query(baseSQL)
	defer res.Close()

	if err != nil {
		return []Post{}, errors.Wrap(err, "cannot get posts"), http.StatusInternalServerError
	}

	post := Post{}

	for res.Next() {
		err := res.Scan(&post.Author, &post.Created, &post.Forum, &post.ID, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)

		if err != nil {
			return []Post{}, errors.Wrap(err, "db query result parsing error"), http.StatusInternalServerError
		}
		sortedPosts = append(sortedPosts, post)
	}

	return sortedPosts, nil, http.StatusOK
}
