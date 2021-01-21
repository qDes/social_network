package controller

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"social_network/internal/config"
	"social_network/internal/model"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var (
	store = sessions.NewCookieStore([]byte("mysession"))
	svc   = config.GetSvc()
)

func Index(resp http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "mysession")
	username := session.Values["username"]
	if username != nil {
		http.Redirect(resp, req, "/account/page/"+fmt.Sprintf("%v", username), http.StatusSeeOther)
	}

	tmp, _ := template.ParseFiles("web/template/index.html")
	tmp.Execute(resp, nil)
}

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
		http.Redirect(resp, req, "/account/page/"+username, http.StatusSeeOther)
	} else {
		data := map[string]interface{}{
			"err": "Invalid",
		}
		tmp, _ := template.ParseFiles("web/template/login/index.html")
		tmp.Execute(resp, data)
	}

}

func UserPage(resp http.ResponseWriter, req *http.Request) {
	var (
		sex string
		add bool
	)
	vars := mux.Vars(req)
	username := vars["username"]

	session, _ := store.Get(req, "mysession")
	sessionUser := session.Values["username"]
	if sessionUser == nil {
		http.Redirect(resp, req, "/", http.StatusSeeOther)
	}

	user, _ := model.GetUser(svc.DB, fmt.Sprintf("%v", username))
	if user.ID == 0 {
		http.Redirect(resp, req, "/", http.StatusSeeOther)
	}
	if bytes.Compare(user.Sex, []byte{1}) == 0 {
		sex = "F"
	} else {
		sex = "M"
	}
	// check add button rendering
	if (username != fmt.Sprintf("%v", sessionUser)) &&
		!(model.CheckFriends(svc.DB, username, fmt.Sprintf("%v", sessionUser))) {
		add = true
	}

	data := map[string]interface{}{
		"username":     username,
		"name":         user.Name,
		"second_name":  user.SecondName,
		"sex":          sex,
		"city":         user.City,
		"interests":    user.Interests,
		"urls":         model.GetFriends(svc.DB, username),
		"add":          add,
		"session_user": sessionUser,
	}
	tmp, _ := template.ParseFiles("web/template/login/user.html")
	tmp.Execute(resp, data)

}

func AddFriend(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	username := req.Form.Get("username")
	sessionUser := req.Form.Get("user_session")
	_ = model.AddFriend(svc.DB, username, sessionUser)
	http.Redirect(resp, req, "/", http.StatusSeeOther)
}

func Logout(resp http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "mysession")
	session.Options.MaxAge = -1
	session.Save(req, resp)
	http.Redirect(resp, req, "/", http.StatusSeeOther)
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
	err := model.InsertUser(svc.DB, username, password, name, secondName, city, interests, sex)
	fmt.Println(err)
	http.Redirect(resp, req, "/account/index", http.StatusSeeOther)
}

func SearchUser(resp http.ResponseWriter, req *http.Request) {
	firstName, ok := req.URL.Query()["firstname"]
	if !ok {
		fmt.Println("Url Param 'firstname' is missing")
	}

	secondName, ok := req.URL.Query()["secondname"]
	if !ok {
		fmt.Println("Url Param 'secondname' is missing")
	}
	fmt.Println(firstName, secondName)
}
