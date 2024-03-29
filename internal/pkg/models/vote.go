package models

import (
	"net/http"

	"github.com/Smet1/tp_db_forum/internal/database"
	"github.com/pkg/errors"
)

//easyjson:json
type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int8   `json:"voice"`
	Thread   int32  `json:"thread"`
}

func CreateVoteAndUpdateThread(voteToCreate Vote) (Thread, error, int) {
	conn := database.Connection

	voiceDiff := voteToCreate.Voice

	resInsert, _ := conn.Exec(`INSERT INTO forum_vote (nickname, voice, thread) VALUES ($1, $2, $3)`,
		voteToCreate.Nickname, voteToCreate.Voice, voteToCreate.Thread)

	if resInsert.RowsAffected() == 0 {

		voteBeforeUpdate, _ := GetVoteByNicknameAndThreadID(voteToCreate.Nickname, voteToCreate.Thread)
		//if err != nil {
		//	return Thread{}, errors.Wrap(err, "Cant find existing vote"), http.StatusInternalServerError
		//}

		voteToCreate, _ := UpdateVote(voteToCreate.Nickname, voteToCreate.Thread, voteToCreate.Voice)
		//if err != nil {
		//	return Thread{}, errors.Wrap(err, "Cant update existing vote"), http.StatusInternalServerError
		//}

		// если меняем отзыв, то нужно откатить предыдущий и накатить новый, поэтому ±2
		if voteToCreate.Voice == -1 && voteToCreate.Voice != voteBeforeUpdate.Voice {
			voiceDiff = -2
		} else if voteToCreate.Voice == 1 && voteToCreate.Voice != voteBeforeUpdate.Voice {
			voiceDiff = 2
		} else if voteToCreate.Voice == voteBeforeUpdate.Voice {
			voiceDiff = 0
		}
	}

	updatedThread, err, status := UpdateThreadVote(voteToCreate.Thread, voiceDiff)
	if err != nil {
		return Thread{}, errors.Wrap(err, "cant update thread"), status
	}

	return updatedThread, nil, http.StatusOK
}

func GetVoteByNicknameAndThreadID(nickname string, threadID int32) (Vote, error) {
	conn := database.Connection

	res, err := conn.Query(`SELECT * FROM forum_vote WHERE nickname = $1 AND thread = $2`, nickname, threadID)

	if err != nil {
		return Vote{}, errors.Wrap(err, "cant find vote")
	}
	defer res.Close()

	existingVote := Vote{}

	if res.Next() {
		err := res.Scan(&existingVote.Nickname, &existingVote.Voice, &existingVote.Thread)
		if err != nil {
			return Vote{}, errors.Wrap(err, "db query result parsing error")
		}

		return existingVote, nil
	}
	return Vote{}, errors.New("cant find vote")
}

func UpdateVote(nickname string, threadID int32, newVoice int8) (Vote, error) {
	conn := database.Connection

	res, err := conn.Exec(`UPDATE forum_vote SET voice = $1 WHERE nickname = $2 AND thread = $3`,
		newVoice, nickname, threadID)
	if err != nil {
		return Vote{}, errors.Wrap(err, "cannot update vote")
	}
	if res.RowsAffected() == 0 {
		return Vote{}, errors.New("not found")
	}

	return Vote{
		Nickname: nickname,
		Voice:    newVoice,
		Thread:   threadID,
	}, nil
}
