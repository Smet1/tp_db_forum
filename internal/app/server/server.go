package server

import (
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"tp_db_forum/internal/pkg/controllers"
)

func Run(port string) error {
	address := ":" + port

	router := mux.NewRouter()

	router.HandleFunc("/api/user/{nickname}/profile", controllers.GetUserProfile).Methods("GET")
	router.HandleFunc("/api/user/{nickname}/profile", controllers.UpdateUserProfile).Methods("POST")
	router.HandleFunc("/api/user/{nickname}/create", controllers.CreateUser).Methods("POST")

	router.HandleFunc("/api/forum/create", controllers.CreateForum).Methods("POST")
	router.HandleFunc("/api/forum/{slug}/details", controllers.GetForum).Methods("GET")
	router.HandleFunc("/api/forum/{slug}/create", controllers.CreateThread).Methods("POST")
	router.HandleFunc("/api/forum/{slug}/threads", controllers.GetThreads).Methods("GET")

	router.HandleFunc("/api/thread/{slug_or_id}/create", controllers.CreatePosts).Methods("POST")
	router.HandleFunc("/api/thread/{slug_or_id}/details", controllers.GetThreadDetails).Methods("GET")

	router.HandleFunc("/api/thread/{slug_or_id}/vote", controllers.CreateVote).Methods("POST")
	router.HandleFunc("/api/thread/{slug_or_id}/posts", controllers.GetThreadPosts).Methods("GET")

	router.HandleFunc("/api/thread/{slug_or_id}/details", controllers.UpdateThread).Methods("POST")

	router.HandleFunc("/api/forum/{slug}/users", controllers.GetForumUsers).Methods("GET")

	router.HandleFunc("/api/service/status", controllers.GetDBStatus).Methods("GET")
	router.HandleFunc("/api/service/clear", controllers.ClearDB).Methods("POST")
	err := http.ListenAndServe(address, router)
	if err != nil {
		return errors.Wrap(err, "server Run error")
	}

	return nil
}
