package server

import (
	"github.com/Smet1/tp_db_forum/internal/pkg/controllers"
	"github.com/pkg/errors"
	goji "goji.io"
	"goji.io/pat"
	"net/http"
)

func Run(port string) error {
	address := ":" + port
	//router := mux.NewRouter()
	//
	//router.HandleFunc("/api/user/{nickname}/profile", controllers.GetUserProfile).Methods("GET")
	//router.HandleFunc("/api/user/{nickname}/profile", controllers.UpdateUserProfile).Methods("POST")
	//router.HandleFunc("/api/user/{nickname}/create", controllers.CreateUser).Methods("POST")
	//
	//router.HandleFunc("/api/forum/create", controllers.CreateForum).Methods("POST")
	//
	//router.HandleFunc("/api/forum/{slug}/details", controllers.GetForum).Methods("GET")
	//router.HandleFunc("/api/forum/{slug}/create", controllers.CreateThread).Methods("POST")
	//router.HandleFunc("/api/forum/{slug}/threads", controllers.GetThreads).Methods("GET")
	//router.HandleFunc("/api/forum/{slug}/users", controllers.GetForumUsers).Methods("GET")
	//
	//router.HandleFunc("/api/thread/{slug_or_id}/create", controllers.CreatePosts).Methods("POST")
	//router.HandleFunc("/api/thread/{slug_or_id}/details", controllers.GetThreadDetails).Methods("GET")
	//router.HandleFunc("/api/thread/{slug_or_id}/details", controllers.UpdateThread).Methods("POST")
	//router.HandleFunc("/api/thread/{slug_or_id}/vote", controllers.CreateVote).Methods("POST")
	//router.HandleFunc("/api/thread/{slug_or_id}/posts", controllers.GetThreadPosts).Methods("GET")
	//
	//router.HandleFunc("/api/service/status", controllers.GetDBStatus).Methods("GET")
	//router.HandleFunc("/api/service/clear", controllers.ClearDB).Methods("POST")
	//
	//router.HandleFunc("/api/post/{id}/details", controllers.UpdatePost).Methods("POST")
	//router.HandleFunc("/api/post/{id}/details", controllers.GetPostInfo).Methods("GET")
	//err := http.ListenAndServe(address, router)
	//if err != nil {
	//	return errors.Wrap(err, "server Run error")
	//}

	gojiRouter := goji.NewMux()
	gojiRouter.HandleFunc(pat.Get("/api/user/:nickname/profile"), controllers.GetUserProfile)
	gojiRouter.HandleFunc(pat.Post("/api/user/:nickname/profile"), controllers.UpdateUserProfile)
	gojiRouter.HandleFunc(pat.Post("/api/user/:nickname/create"), controllers.CreateUser)

	gojiRouter.HandleFunc(pat.Post("/api/forum/create"), controllers.CreateForum)

	gojiRouter.HandleFunc(pat.Get("/api/forum/{slug}/details"), controllers.GetForum)
	gojiRouter.HandleFunc(pat.Post("/api/forum/{slug}/create"), controllers.CreateThread)
	gojiRouter.HandleFunc(pat.Get("/api/forum/{slug}/threads"), controllers.GetThreads)
	gojiRouter.HandleFunc(pat.Get("/api/forum/{slug}/users"), controllers.GetForumUsers)

	gojiRouter.HandleFunc(pat.Post("/api/thread/{slug_or_id}/create"), controllers.CreatePosts)
	gojiRouter.HandleFunc(pat.Get("/api/thread/{slug_or_id}/details"), controllers.GetThreadDetails)
	gojiRouter.HandleFunc(pat.Post("/api/thread/{slug_or_id}/details"), controllers.UpdateThread)
	gojiRouter.HandleFunc(pat.Post("/api/thread/{slug_or_id}/vote"), controllers.CreateVote)
	gojiRouter.HandleFunc(pat.Get("/api/thread/{slug_or_id}/posts"), controllers.GetThreadPosts)

	gojiRouter.HandleFunc(pat.Get("/api/service/status"), controllers.GetDBStatus)
	gojiRouter.HandleFunc(pat.Post("/api/service/clear"), controllers.ClearDB)

	gojiRouter.HandleFunc(pat.Post("/api/post/{id}/details"), controllers.UpdatePost)
	gojiRouter.HandleFunc(pat.Get("/api/post/{id}/details"), controllers.GetPostInfo)

	err := http.ListenAndServe(address, gojiRouter)
	if err != nil {
		return errors.Wrap(err, "server Run error")
	}

	return nil
}
