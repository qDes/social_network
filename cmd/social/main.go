package main

import (
	"log"
	"net/http"
	"social_network/internal/controller"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	//
	go controller.InitFeedCache()
	//read queue
	go controller.QReader()

	r := mux.NewRouter()

	r.HandleFunc("/", controller.Index)
	r.HandleFunc("/account/index", controller.IndexLogin)
	r.HandleFunc("/account/login", controller.Login)
	r.HandleFunc("/account/page/{username}", controller.UserPage)
	r.HandleFunc("/account/logout", controller.Logout)
	r.HandleFunc("/account/signupindex", controller.SignUpIndex)
	r.HandleFunc("/account/signup", controller.SignUp)
	r.HandleFunc("/account/add_friend", controller.AddFriend)
	r.HandleFunc("/account/search_user", controller.SearchUser)
	r.HandleFunc("/account/search_user_tarantool", controller.SearchUserT)
	r.HandleFunc("/account/search", controller.Search)

	r.HandleFunc("/account/feed", controller.UserFeed)
	r.HandleFunc("/account/add_post", controller.AddUserPost)


	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:3000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
