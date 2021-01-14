package controller

import (
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("mysession"))

func Index(resp http.ResponseWriter, req *http.Request) {
	tmp, _ := template.ParseFiles("web/template/index.html")
	tmp.Execute(resp, nil)
}

func Login(resp http.ResponseWriter, req *http.Request) {

	req.ParseForm()
	username := req.Form.Get("username")
	password := req.Form.Get("password")

	if username == "abc" && password == "123" {
		session, _ := store.Get(req, "mysession")
		session.Values["username"] = username
		session.Save(req, resp)
		http.Redirect(resp, req, "/account/welcome", http.StatusSeeOther)
	} else {
		data := map[string]interface{}{
			"err": "Invalid",
		}
		tmp, _ := template.ParseFiles("web/template/index.html")
		tmp.Execute(resp, data)
	}

}

func Welcome(resp http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "mysession")
	username := session.Values["username"]
	data := map[string]interface{}{
		"username": username,
	}
	tmp, _ := template.ParseFiles("web/template/welcome.html")
	tmp.Execute(resp, data)

}

func Logout(resp http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "mysession")
	session.Options.MaxAge = -1
	session.Save(req, resp)
	http.Redirect(resp, req, "/account/index", http.StatusSeeOther)
}
