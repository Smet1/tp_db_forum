package models

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	Path     []int64   `json:"-"`
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
	// TODO(): транзакция, по тз не пройдет (при косяке в 1 посте остальные до него заносятся в БД)
	conn := database.Connection

	now := time.Now()
	var id int64 = 0
	// get last id
	res, err := conn.Query(`SELECT last_value FROM forum_post_id_seq`)
	if err != nil {
		return []Post{}, errors.Wrap(err, "cant get last id"), http.StatusInternalServerError
	}
	for res.Next() {
		err := res.Scan(&id)

		if err != nil {
			return []Post{}, errors.Wrap(err, "db query result parsing error"), http.StatusInternalServerError
		}
	}
	res.Close()
	id++
	log.Println("\tlast id =", id)

	for i, post := range postsToCreate {
		fmt.Println("--", i, "--")

		if post.Parent != 0 {
			parentPostQuery, err, _ := GetPostByID(post.Parent)
			if err != nil {
				return []Post{}, errors.Wrap(err, "cant get parent post"), http.StatusConflict
			}

			fmt.Println("--== check thread id ==--")
			fmt.Println("\tcurrent thread =", existingThread.ID)
			fmt.Println("\tthread in post founded =", parentPostQuery.Thread)
			fmt.Println("\tparent in post =", post.Parent)
			fmt.Println("--==                 ==--")

			if parentPostQuery.Thread != existingThread.ID {
				return []Post{}, errors.New("parent post created in another thread"), http.StatusConflict
			}

			post.Path = parentPostQuery.Path
		}

		post.Path = append(post.Path, id)

		resInsert, err := conn.Exec(`INSERT INTO forum_post (author, created, forum, message, parent, thread, path) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			post.Author, now, existingThread.Forum, post.Message, post.Parent, existingThread.ID,
			"{"+strings.Trim(strings.Replace(fmt.Sprint(post.Path), " ", ",", -1), "[]")+"}")

		fmt.Println("\tpath = ", strings.Trim(strings.Replace(fmt.Sprint(post.Path), " ", ",", -1), "[]"))

		//if err != nil {
		//	return []Post{}, errors.Wrap(err, "cant insert post"), http.StatusInternalServerError
		//}

		if resInsert.RowsAffected() == 0 {
			return []Post{}, errors.Wrap(err, "cant create thread"), http.StatusNotFound
		}

		postsToCreate[i].Forum = existingThread.Forum
		postsToCreate[i].Thread = existingThread.ID
		postsToCreate[i].Created = now
		postsToCreate[i].ID = id
		id++

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
	if sort == "" {
		sort = "flat"
	}
	// tree && parent tree UNDONE
	conn := database.Connection

	baseSQL := ""
	sortedPosts := make([]Post, 0, 1)
	//strID := strconv.FormatInt(int64(parentThread.ID), 10)

	switch sort {
	case "flat":
		baseSQL = FlatSort(parentThread, limit, since, sort, desc)

		fmt.Println("---===flat sort===---")
		fmt.Println("\tbaseSQL =", baseSQL)

	case "tree":
		baseSQL = TreeSort(parentThread, limit, since, sort, desc)

		fmt.Println("---===tree sort===---")
		fmt.Println("\tbaseSQL =", baseSQL)

	case "parent_tree":
		rootPosts, err := ParentTreeSort(parentThread, limit, since, sort, desc)
		if err != nil {
			return []Post{}, errors.Wrap(err, "cannot get posts"), http.StatusInternalServerError
		}

		if len(rootPosts) == 0 {
			log.Println("parent tree: no posts found")

			return []Post{}, nil, http.StatusOK
		}

		for _, val := range rootPosts {
			//sortedPosts = append(sortedPosts, val)

			childPostsQuery, err := conn.Query("SELECT author, created, forum, id, isedited, message, parent, thread" +
				" FROM forum_post" +
				" WHERE path[1] = $1 ORDER BY path", val.ID)

			if err != nil {
				childPostsQuery.Close()

				return []Post{}, errors.Wrap(err, "db query result parsing error"), http.StatusInternalServerError
			}

			post := Post{}

			for childPostsQuery.Next() {
				err := childPostsQuery.Scan(&post.Author, &post.Created, &post.Forum, &post.ID, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)

				if err != nil {
					childPostsQuery.Close()

					return []Post{}, errors.Wrap(err, "db query result parsing error"), http.StatusInternalServerError
				}
				sortedPosts = append(sortedPosts, post)
			}

			childPostsQuery.Close()
		}

		return sortedPosts, nil, http.StatusOK
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

	if len(sortedPosts) == 0 {
		log.Println("GetSortedPosts: (not parent_tree sort) not posts found")

		return []Post{}, nil, http.StatusOK
	}

	return sortedPosts, nil, http.StatusOK
}

func FlatSort(parentThread Thread, limit int, since int, sort string, desc bool) string {
	strID := strconv.FormatInt(int64(parentThread.ID), 10)
	baseSQL := ""

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

	return baseSQL
}

func TreeSort(parentThread Thread, limit int, since int, sort string, desc bool) string {
	strID := strconv.FormatInt(int64(parentThread.ID), 10)
	baseSQL := ""

	baseSQL = "SELECT author, created, forum, id, isedited, message, parent, thread FROM forum_post WHERE thread = " + strID

	if since != 0 {
		if desc {
			baseSQL += " AND path < (SELECT path FROM forum_post WHERE id = " + strconv.Itoa(since) + ")"
		} else {
			baseSQL += " AND path > (SELECT path FROM forum_post WHERE id = " + strconv.Itoa(since) + ")"
		}
	}

	if desc {
		baseSQL += " ORDER BY path DESC, id DESC"
	} else {
		baseSQL += " ORDER BY path, id"
	}

	baseSQL += " LIMIT " + strconv.Itoa(limit)

	return baseSQL
}

func ParentTreeSort(parentThread Thread, limit int, since int, sort string, desc bool) ([]Post, error) {
	conn := database.Connection
	strID := strconv.FormatInt(int64(parentThread.ID), 10)
	baseSQL := ""

	baseSQL = "SELECT author, created, forum, id, isedited, message, parent, thread FROM forum_post WHERE thread = " + strID + " AND parent = 0"

	if since != 0 {
		if desc {
			baseSQL += " AND path[1] < (SELECT path[1] FROM forum_post WHERE id = " + strconv.Itoa(since) + ")"
		} else {
			baseSQL += " AND path[1] > (SELECT path[1] FROM forum_post WHERE id = " + strconv.Itoa(since) + ")"
		}
	}

	if desc {
		baseSQL += " ORDER BY id DESC"
	} else {
		baseSQL += " ORDER BY id"
	}

	baseSQL += " LIMIT " + strconv.Itoa(limit)

	fmt.Println("---===parent_tree sort===---")
	fmt.Println("\tbaseSQL =", baseSQL)

	rootPostsRaw, err := conn.Query(baseSQL)
	defer rootPostsRaw.Close()

	if err != nil {
		return []Post{}, errors.Wrap(err, "db query result parsing error")
	}

	post := Post{}
	rootPosts := make([]Post, 0, 1)

	for rootPostsRaw.Next() {
		err := rootPostsRaw.Scan(&post.Author, &post.Created, &post.Forum, &post.ID, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)

		if err != nil {
			return []Post{}, errors.Wrap(err, "db query result parsing error")
		}
		rootPosts = append(rootPosts, post)
	}

	return rootPosts, nil
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

	res, err := conn.Query("SELECT author, created, forum, id, isedited, message, parent, thread, path FROM forum_post WHERE id = $1", id)
	defer res.Close()

	if err != nil {
		return Post{}, errors.Wrap(err, "cannot get post"), http.StatusNotFound
	}

	post := Post{}

	for res.Next() {
		//path := []pgtype.ArrayDimension{}
		err := res.Scan(&post.Author, &post.Created, &post.Forum, &post.ID, &post.IsEdited, &post.Message, &post.Parent, &post.Thread, pq.Array(&post.Path))

		if err != nil {
			return Post{}, errors.Wrap(err, "db query result parsing error"), http.StatusInternalServerError
		}

		fmt.Println("GetPostByID::path = ", post.Path)

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
