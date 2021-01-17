package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"social_network/internal/config"
	"social_network/internal/model"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var (
	store = sessions.NewCookieStore([]byte("mysession"))
	svc   = config.GetDB()
)

func IndexLogin(resp http.ResponseWriter, req *http.Request) {
	tmp, _ := template.ParseFiles("web/template/login/index.html")
	tmp.Execute(resp, nil)
}

func Login(resp http.ResponseWriter, req *http.Request) {

	req.ParseForm()
	username := req.Form.Get("username")
	password := req.Form.Get("password")

	dbPass, _ := model.GetPass(svc.DB, username)
	err := bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(password))
	if err == nil {
		session, _ := store.Get(req, "mysession")
		session.Values["username"] = username
		session.Save(req, resp)
		http.Redirect(resp, req, "/account/welcome", http.StatusSeeOther)
	} else {
		data := map[string]interface{}{
			"err": "Invalid",
		}
		tmp, _ := template.ParseFiles("web/template/login/index.html")
		tmp.Execute(resp, data)
	}

}

func UserPage(resp http.ResponseWriter, req *http.Request) {
	var sex string
	session, _ := store.Get(req, "mysession")
	username := session.Values["username"]
	user, _ := model.GetUser(svc.DB, fmt.Sprintf("%v", username))
	if user.Sex {
		sex = "F"
	} else {
		sex = "M"
	}
	data := map[string]interface{}{
		"username": username,
		"name": user.Name,
		"second_name": user.SecondName,
		"sex": sex,
		"city": user.City,
		"interests": user.Interests,
	}
	tmp, _ := template.ParseFiles("web/template/login/user.html")
	tmp.Execute(resp, data)

}

func Logout(resp http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "mysession")
	session.Options.MaxAge = -1
	session.Save(req, resp)
	http.Redirect(resp, req, "/account/index", http.StatusSeeOther)
}


func SignUpIndex(resp http.ResponseWriter, req *http.Request) {
	tmp, _ := template.ParseFiles("web/template/signup/index.html")
	tmp.Execute(resp, nil)
}

func SignUp(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	// TODO: check unique username or db constrains
	username := req.Form.Get("username")
	password := getHash([]byte(req.Form.Get("password")))
	name := req.Form.Get("name")
	secondName := req.Form.Get("second_name")
	sexIn := req.Form.Get("sex")
	city := req.Form.Get("city")
	interests := req.Form.Get("interests")
	var sex bool
	if sexIn == "0" {
		sex = false
	} else {
		sex = true
	}
	//fmt.Println(username, password, name, secondName, sex, city, interests)
	err := model.InsertUser(svc.DB, username, password, name, secondName, city, interests, sex)
	fmt.Println(err)
	http.Redirect(resp, req, "/account/index", http.StatusSeeOther)
}
