package models

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
	"time"
	"tp_db_forum/internal/database"
)

type Thread struct {
	Author  string    `json:"author"`
	Created time.Time `json:"created"`
	Forum   string    `json:"forum"`
	ID      int32     `json:"id"`
	Message string    `json:"message"`
	Slug    string    `json:"slug"`
	Title   string    `json:"title"`
	Votes   int32     `json:"votes"`
}

func CreateThread(threadToCreate Thread) (Thread, error, int) {
	conn := database.Connection

	existingForum, err := GetForumBySlug(threadToCreate.Forum)
	if err != nil {
		return Thread{}, errors.Wrap(err, "cant find slug"), http.StatusNotFound
	}
	threadToCreate.Forum = existingForum.Slug
	fmt.Println("\t\tDB forum.slug = ", existingForum.Slug)

	existingUser, err := GetUserByNickname(threadToCreate.Author)
	if err != nil {
		return Thread{}, errors.Wrap(err, "cant find user"), http.StatusNotFound
	}
	threadToCreate.Author = existingUser.Nickname
	fmt.Println("\t\tDB user.nickname = ", existingUser.Nickname)

	if threadToCreate.Slug == "" {
		resInsert, err := conn.Exec(`INSERT INTO forum_thread (author, created, forum, message, slug, title) VALUES ($1, $2, $3, $4, NULL, $5)`,
			threadToCreate.Author, threadToCreate.Created, threadToCreate.Forum, threadToCreate.Message, threadToCreate.Title)
		if resInsert.RowsAffected() == 0 {
			return Thread{}, errors.Wrap(err, "cant create thread"), http.StatusConflict
		}
	} else {
		resInsert, err := conn.Exec(`INSERT INTO forum_thread (author, created, forum, message, slug, title) VALUES ($1, $2, $3, $4, $5, $6)`,
			threadToCreate.Author, threadToCreate.Created, threadToCreate.Forum, threadToCreate.Message, threadToCreate.Slug,
			threadToCreate.Title)
		if resInsert.RowsAffected() == 0 {
			return Thread{}, errors.Wrap(err, "cant create thread"), http.StatusConflict
		}
	}

	//resInsert, err := conn.Exec(`INSERT INTO forum_thread (author, created, forum, message, slug, title) VALUES ($1, $2, $3, $4, $5, $6)`,
	//	threadToCreate.Author, threadToCreate.Created, threadToCreate.Forum, threadToCreate.Message, threadToCreate.Slug,
	//	threadToCreate.Title)
	//if resInsert.RowsAffected() == 0 {
	//	return Thread{}, errors.Wrap(err, "cant create thread"), http.StatusInternalServerError
	//}

	// get last id
	res, err := conn.Query(`SELECT last_value FROM forum_thread_id_seq`)
	for res.Next() {
		err := res.Scan(&threadToCreate.ID)

		if err != nil {
			return Thread{}, errors.Wrap(err, "db query result parsing error"), http.StatusInternalServerError
		}
	}
	log.Println("\t\t CreateThread id = ", threadToCreate.ID)

	//res, err := conn.Query(`SELECT id FROM forum_thread WHERE slug = $1`, threadToCreate.Slug)
	//defer res.Close()
	//
	//for res.Next() {
	//	err := res.Scan(&threadToCreate.ID)
	//
	//	if err != nil {
	//		return Thread{}, errors.Wrap(err, "db query result parsing error"), http.StatusInternalServerError
	//	}
	//}
	//log.Println("\t\t CreateThread id = ", threadToCreate.ID)

	status := UpdateForumStats(existingForum, "thread", true, 1)
	if status != http.StatusOK {
		return Thread{}, errors.New("cant update forum stats"), status
	}

	return threadToCreate, nil, http.StatusOK
}

func GetForumThreads(slug string, limit int, since string, desc bool) ([]Thread, error, int) {
	conn := database.Connection

	queriedThreads := make([]Thread, 0, 1)

	existingForum, err := GetForumBySlug(slug)
	if err != nil {
		return []Thread{}, errors.Wrap(err, "cant find slug"), http.StatusNotFound
	}

	baseSQL := "SELECT * FROM forum_thread"

	baseSQL += " WHERE forum = '" + existingForum.Slug + "'"

	if since != "" {
		if desc {
			baseSQL += " AND created <= '" + since + "'" // ::timestamptz
		} else {
			baseSQL += " AND created >= '" + since + "'"
		}
	}

	if desc {
		baseSQL += " ORDER BY created DESC"
	} else {
		baseSQL += " ORDER BY created"
	}

	if limit > 0 {
		baseSQL += " LIMIT " + strconv.Itoa(limit)
	}

	log.Println(baseSQL)
	res, err := conn.Query(baseSQL)
	defer res.Close()

	if err != nil {
		return []Thread{}, errors.Wrap(err, "cannot get user by nickname or email"), http.StatusInternalServerError
	}

	t := Thread{}

	for res.Next() {
		err := res.Scan(&t.Author, &t.Created, &t.Forum, &t.ID, &t.Message, &t.Slug, &t.Title, &t.Votes)

		if err != nil {
			return []Thread{}, errors.Wrap(err, "db query result parsing error"), http.StatusInternalServerError
		}
		queriedThreads = append(queriedThreads, t)
	}

	return queriedThreads, nil, http.StatusOK
}

func GetThreadByIDorSlug(id int, slug string) (Thread, error, int) {
	conn := database.Connection
	t := Thread{}
	if id == -1 {

		res, err := conn.Query(`SELECT * FROM forum_thread WHERE slug = $1`, slug)
		defer res.Close()

		if err != nil {
			return Thread{}, errors.Wrap(err, "cannot get thread by slug"), http.StatusNotFound
		}

		if res.Next() {
			err := res.Scan(&t.Author, &t.Created, &t.Forum, &t.ID, &t.Message, &t.Slug, &t.Title, &t.Votes)
			if err != nil {
				return Thread{}, errors.Wrap(err, "db query result parsing error"), http.StatusInternalServerError
			}

			return t, nil, http.StatusOK
		}

		return Thread{}, errors.New("cannot get thread by nickname"), http.StatusNotFound

	} else if slug == "" {

		res, err := conn.Query(`SELECT * FROM forum_thread WHERE id = $1`, id)
		defer res.Close()

		if err != nil {
			return Thread{}, errors.Wrap(err, "cannot get thread by id"), http.StatusNotFound
		}

		if res.Next() {
			err := res.Scan(&t.Author, &t.Created, &t.Forum, &t.ID, &t.Message, &t.Slug, &t.Title, &t.Votes)
			if err != nil {
				return Thread{}, errors.Wrap(err, "db query result parsing error"), http.StatusInternalServerError
			}

			return t, nil, http.StatusOK
		}

		return Thread{}, errors.New("cannot get thread by id"), http.StatusNotFound

	} else {
		res, err := conn.Query(`SELECT * FROM forum_thread WHERE id = $1 OR slug = $2`, id, slug)
		defer res.Close()

		if err != nil {
			return Thread{}, errors.Wrap(err, "cannot get thread by id or slug"), http.StatusNotFound
		}

		if res.Next() {
			err := res.Scan(&t.Author, &t.Created, &t.Forum, &t.ID, &t.Message, &t.Slug, &t.Title, &t.Votes)
			if err != nil {
				return Thread{}, errors.Wrap(err, "db query result parsing error"), http.StatusInternalServerError
			}

			return t, nil, http.StatusOK
		}

		return Thread{}, errors.New("cannot get thread by id or slug"), http.StatusNotFound
	}
}

func UpdateThreadVote(threadId int32, voteValue int8) (Thread, error, int) {
	conn := database.Connection

	fmt.Println("UpdateThreadVote")
	fmt.Println(voteValue)
	plus := "+"
	if voteValue < 0 {
		plus = "-"
		voteValue = -voteValue
	}

	fmt.Println(plus)

	res, err := conn.Exec(`UPDATE forum_thread SET votes = votes`+plus+`$1 WHERE id = $2`, voteValue, threadId)
	if err != nil {
		return Thread{}, errors.Wrap(err, "cannot update thread"), http.StatusConflict
	}

	if res.RowsAffected() == 0 {
		return Thread{}, errors.New("not found"), http.StatusNotFound
	}

	updatedThread, _, _ := GetThreadByIDorSlug(int(threadId), "")
	return updatedThread, nil, http.StatusOK
}

func UpdateThread(existingThread Thread, newThread Thread) (Thread, error, int) {
	conn := database.Connection

	if newThread.Message == "" && newThread.Title == "" {
		return existingThread, nil, http.StatusOK
	}

	baseSQL := "UPDATE forum_thread SET"
	if newThread.Message == "" {
		baseSQL += " message = message,"
	} else {
		baseSQL += " message = '" + newThread.Message + "',"
	}

	if newThread.Title == "" {
		baseSQL += " title = title"
	} else {
		baseSQL += " title = '" + newThread.Title + "'"
	}

	baseSQL += " WHERE slug = '" + existingThread.Slug + "'"

	fmt.Println("\t", baseSQL)
	res, err := conn.Exec(baseSQL)
	if err != nil {
		return Thread{}, errors.Wrap(err, "cannot update thread"), http.StatusConflict
	}

	if res.RowsAffected() == 0 {
		return Thread{}, errors.New("not found"), http.StatusNotFound
	}

	updatedThread, _, _ := GetThreadByIDorSlug(-1, existingThread.Slug)

	return updatedThread, nil, http.StatusOK
}
