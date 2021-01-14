package main

import (
	"net/http"
	"social_network/internal/controller"
)

func main() {
	http.HandleFunc("/account", controller.Index)
	http.HandleFunc("/account/index", controller.Index)
	http.HandleFunc("/account/login", controller.Login)
	http.HandleFunc("/account/welcome", controller.Welcome)
	http.HandleFunc("/account/logout", controller.Logout)

	http.ListenAndServe(":3000", nil)
}
