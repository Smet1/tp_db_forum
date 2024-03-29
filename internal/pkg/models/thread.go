package models

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Smet1/tp_db_forum/internal/database"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"github.com/pkg/errors"
)

//easyjson:json
type Thread struct {
	Author  string    `json:"author"`
	Created time.Time `json:"created,omitempty"`
	Forum   string    `json:"forum,omitempty"`
	ID      int32     `json:"id,omitempty"`
	Message string    `json:"message"`
	Slug    string    `json:"slug,omitempty"`
	Title   string    `json:"title"`
	Votes   int32     `json:"votes,omitempty"`
}

//easyjson:json
type Threads []Thread

func CreateThread(threadToCreate Thread) (Thread, error, int) {
	conn := database.Connection

	existingForum, err := GetForumBySlug(threadToCreate.Forum)
	if err != nil {
		return Thread{}, errors.Wrap(err, "cant find slug"), http.StatusNotFound
	}
	threadToCreate.Forum = existingForum.Slug

	existingUser, err := GetUserByNickname(threadToCreate.Author)
	if err != nil {
		return Thread{}, errors.Wrap(err, "cant find user"), http.StatusNotFound
	}
	threadToCreate.Author = existingUser.Nickname

	tx, _ := conn.Begin()
	defer tx.Rollback()

	if threadToCreate.Slug == "" {
		err := tx.QueryRow(`INSERT INTO forum_thread (author, created, forum, message, slug, title) VALUES ($1, $2, $3, $4, NULL, $5) RETURNING id`,
			threadToCreate.Author, threadToCreate.Created, threadToCreate.Forum, threadToCreate.Message,
			threadToCreate.Title).Scan(&threadToCreate.ID)

		if err == pgx.ErrNoRows {
			return Thread{}, errors.Wrap(err, "cant create thread"), http.StatusConflict
		} else if err != nil {
			return Thread{}, errors.Wrap(err, "cant create thread"), http.StatusConflict
		}
		//if resInsert.RowsAffected() == 0 {
		//	return Thread{}, errors.Wrap(err, "cant create thread"), http.StatusConflict
		//}
	} else {
		err := tx.QueryRow(`INSERT INTO forum_thread (author, created, forum, message, slug, title) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
			threadToCreate.Author, threadToCreate.Created, threadToCreate.Forum, threadToCreate.Message, threadToCreate.Slug,
			threadToCreate.Title).Scan(&threadToCreate.ID)

		if err == pgx.ErrNoRows {
			return Thread{}, errors.Wrap(err, "cant create thread"), http.StatusConflict
		} else if err != nil {
			return Thread{}, errors.Wrap(err, "cant create thread"), http.StatusConflict
		}
		//if resInsert.RowsAffected() == 0 {
		//	return Thread{}, errors.Wrap(err, "cant create thread"), http.StatusConflict
		//}
	}

	tx.Commit()

	status := UpdateForumStats(existingForum, "thread", true, 1)
	if status != http.StatusOK {
		return Thread{}, errors.New("cant update forum stats"), status
	}

	AddUser(threadToCreate.Author, existingForum.Slug)

	return threadToCreate, nil, http.StatusOK
}

func GetForumThreads(slug string, limit int, since string, desc bool) (Threads, error, int) {
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
			baseSQL += " AND created <= '" + since + "'"
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

	res, _ := conn.Query(baseSQL)
	//if err != nil {
	//	return []Thread{}, errors.Wrap(err, "cannot get user by nickname or email"), http.StatusInternalServerError
	//}
	defer res.Close()

	t := Thread{}
	nullSlug := &pgtype.Varchar{}
	for res.Next() {
		_ = res.Scan(&t.Author, &t.Created, &t.Forum, &t.ID, &t.Message, nullSlug, &t.Title, &t.Votes)

		//if err != nil {
		//	return []Thread{}, errors.Wrap(err, "db query result parsing error"), http.StatusInternalServerError
		//}
		t.Slug = nullSlug.String

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
			nullString := pgtype.Text{}
			_ = res.Scan(&t.Author, &t.Created, &t.Forum, &t.ID, &t.Message, &nullString, &t.Title, &t.Votes)
			//if err != nil {
			//	return Thread{}, errors.Wrap(err, "db query result parsing error"), http.StatusInternalServerError
			//}

			t.Slug = nullString.String

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
			nullString := pgtype.Text{}
			_ = res.Scan(&t.Author, &t.Created, &t.Forum, &t.ID, &t.Message, &nullString, &t.Title, &t.Votes)
			//if err != nil {
			//	return Thread{}, errors.Wrap(err, "db query result parsing error"), http.StatusInternalServerError
			//}

			t.Slug = nullString.String

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
			nullString := pgtype.Text{}
			_ = res.Scan(&t.Author, &t.Created, &t.Forum, &t.ID, &t.Message, &nullString, &t.Title, &t.Votes)
			//if err != nil {
			//	return Thread{}, errors.Wrap(err, "db query result parsing error"), http.StatusInternalServerError
			//}

			t.Slug = nullString.String

			return t, nil, http.StatusOK
		}

		return Thread{}, errors.New("cannot get thread by id or slug"), http.StatusNotFound
	}
}

func UpdateThreadVote(threadId int32, voteValue int8) (Thread, error, int) {
	conn := database.Connection
	tx, _ := conn.Begin()
	defer tx.Rollback()

	updatedThread := Thread{}
	slugNullable := &pgtype.Varchar{}
	err := tx.QueryRow(`UPDATE forum_thread SET votes = votes+$1 WHERE id = $2 RETURNING author, created, forum, "message", slug, title, id, votes`,
		voteValue, threadId).Scan(&updatedThread.Author, &updatedThread.Created, &updatedThread.Forum,
		&updatedThread.Message, slugNullable, &updatedThread.Title, &updatedThread.ID, &updatedThread.Votes)
	updatedThread.Slug = slugNullable.String
	if err != nil {
		return Thread{}, errors.New("not found"), http.StatusNotFound
	}
	tx.Commit()

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
