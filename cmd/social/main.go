package main

import (
	"net/http"
	"social_network/internal/controller"
)

func main() {
	http.HandleFunc("/account", controller.IndexLogin)
	http.HandleFunc("/account/index", controller.IndexLogin)
	http.HandleFunc("/account/login", controller.Login)
	http.HandleFunc("/account/welcome", controller.UserPage)
	http.HandleFunc("/account/logout", controller.Logout)
	http.HandleFunc("/account/signupindex", controller.SignUpIndex)
	http.HandleFunc("/account/signup", controller.SignUp)

	http.ListenAndServe(":3000", nil)
}
