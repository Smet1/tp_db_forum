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

type PostFull struct {
	Author *User   `json:"author,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
	Post   *Post   `json:"post,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
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

		if post.Parent != 0 {
			parentPostQuery, err, _ := GetPostByID(post.Parent)
			if err != nil {
				return []Post{}, errors.Wrap(err, "cant get parent post"), http.StatusConflict
			}

			if parentPostQuery.Thread != existingThread.ID {
				return []Post{}, errors.New("parent post created in another thread"), http.StatusConflict
			}
		}

		resInsert, err := conn.Exec(`INSERT INTO forum_post (author, created, forum, message, parent, thread) VALUES ($1, $2, $3, $4, $5, $6)`,
			post.Author, now, existingThread.Forum, post.Message, post.Parent, existingThread.ID)

		if resInsert.RowsAffected() == 0 {
			return []Post{}, errors.Wrap(err, "cant create thread"), http.StatusNotFound
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

	existingForum, _ := GetForumBySlug(existingThread.Forum)
	status := UpdateForumStats(existingForum, "post", true, len(postsToCreate))
	if status != http.StatusOK {
		return []Post{}, errors.New("cant update forum stats"), status
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

func UpdatePost(existingPost Post, newPost Post) (Post, error, int) {
	conn := database.Connection

	if newPost.Message == "" {
		return existingPost, nil, http.StatusOK
	}

	if existingPost.Message == newPost.Message {
		return existingPost, nil, http.StatusOK
	}

	res, err := conn.Exec("UPDATE forum_post SET message = $1, isedited = true WHERE id = $2", newPost.Message, existingPost.ID)
	if err != nil {
		return Post{}, errors.Wrap(err, "cannot update post"), http.StatusConflict
	}

	if res.RowsAffected() == 0 {
		return Post{}, errors.New("not found"), http.StatusNotFound
	}

	updatedPost, _, _ := GetPostByID(existingPost.ID)

	return updatedPost, nil, http.StatusOK
}

func GetPostByID(id int64) (Post, error, int) {
	conn := database.Connection

	res, err := conn.Query("SELECT author, created, forum, id, isedited, message, parent, thread FROM forum_post WHERE id = $1", id)
	defer res.Close()

	if err != nil {
		return Post{}, errors.Wrap(err, "cannot get post"), http.StatusNotFound
	}

	post := Post{}

	for res.Next() {
		err := res.Scan(&post.Author, &post.Created, &post.Forum, &post.ID, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)

		if err != nil {
			return Post{}, errors.Wrap(err, "db query result parsing error"), http.StatusInternalServerError
		}

		return post, nil, http.StatusOK
	}

	return Post{}, errors.New("cant find post with this id"), http.StatusNotFound
}

func GetPostDetails(existingPost Post, related []string) (PostFull, error, int) {
	conn := database.Connection
	baseSQL := ""
	postFull := PostFull{}
	for _, val := range related {

		switch val {
		case "user":
			baseSQL = `SELECT about, email, fullname, nickname FROM forum_users WHERE nickname = $1`
			res, _ := conn.Query(baseSQL, existingPost.Author)

			u := User{}

			for res.Next() {
				_ = res.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname)
			}

			postFull.Author = &u

		case "forum":
			baseSQL = `SELECT posts, slug, threads, title, "user" FROM forum_forum WHERE slug = $1`
			res, _ := conn.Query(baseSQL, existingPost.Forum)

			f := Forum{}

			for res.Next() {
				_ = res.Scan(&f.Posts, &f.Slug, &f.Threads, &f.Title, &f.User)
			}

			postFull.Forum = &f
		case "thread":
			baseSQL = `SELECT author, created, forum, id, message, slug, title, votes FROM forum_thread WHERE id = $1`
			res, _ := conn.Query(baseSQL, existingPost.Thread)

			t := Thread{}

			for res.Next() {
				_ = res.Scan(&t.Author, &t.Created, &t.Forum, &t.ID, &t.Message, &t.Slug, &t.Title, &t.Votes)
			}

			postFull.Thread = &t
		}
	}

	return postFull, nil, http.StatusOK
}
