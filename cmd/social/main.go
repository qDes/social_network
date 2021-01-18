package main

import (
	"log"
	"net/http"
	"social_network/internal/controller"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", controller.Index)
	r.HandleFunc("/account/index", controller.IndexLogin)
	r.HandleFunc("/account/login", controller.Login)
	r.HandleFunc("/account/page/{username}", controller.UserPage)
	r.HandleFunc("/account/logout", controller.Logout)
	r.HandleFunc("/account/signupindex", controller.SignUpIndex)
	r.HandleFunc("/account/signup", controller.SignUp)
	r.HandleFunc("/account/add_friend", controller.AddFriend)

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:3000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
