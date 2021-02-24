package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"social_network/internal/config"
	"social_network/internal/model"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/streadway/amqp"
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

	//tmp, _ := template.ParseFiles("web/template/test/base.html", "web/template/test/content.html")
	//tmp.Execute(resp, nil)
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

func UserFeed(resp http.ResponseWriter, req *http.Request) {
	var feed model.Feed


	session, _ := store.Get(req, "mysession")
	sessionUser := session.Values["username"]
	if sessionUser == nil {
		http.Redirect(resp, req, "/", http.StatusSeeOther)
	}
	user, _ := model.GetUser(svc.DB, fmt.Sprintf("%v", sessionUser))
	if user.ID == 0 {
		http.Redirect(resp, req, "/", http.StatusSeeOther)
	}

	ctx := context.Background()
	redisFeed := svc.RDB.Get(ctx, strconv.Itoa(int(user.ID)))
	redisBytes, err := redisFeed.Bytes()
	if err != nil {
		fmt.Println("redis Bytes error")
	}
	err = json.Unmarshal(redisBytes, &feed)
	if err != nil {
		fmt.Println("unmarshalling error")
	}

	data := map[string]interface{}{
		"posts": feed.Posts,
	}

	tmp, _ := template.ParseFiles("web/template/feed/feed.html")
	tmp.Execute(resp, data)
}

func AddUserPost(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	userID, _ := strconv.ParseInt(req.Form.Get("user_id"), 10, 64)
	text := req.Form.Get("post_text")
	//write post to mysql
	now := time.Now()
	err := model.SavePost(svc.DB, userID, text, now)
	if err != nil {
		fmt.Println(err)
	}
	// get user friends
	friendsIDs := model.GetFriendsIDs(svc.DB, userID)
	// write to rabbit friend_id + user post

	for _, friendID := range friendsIDs {
		msg := model.Post{FriendID: friendID, UserID: userID, Text: text, Date: now}
		data, err := json.Marshal(msg)
		err = svc.Feed.Publish(
			"",         // exchange
			svc.Q.Name, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        data,
			})
		if err != nil {
			fmt.Println(err)
		}
	}

	http.Redirect(resp, req, "/", http.StatusSeeOther)

}

func UserPage(resp http.ResponseWriter, req *http.Request) {
	var (
		sex          string
		add, addPost bool
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

	//check add post
	if username == fmt.Sprintf("%v", sessionUser) {
		addPost = true
	}

	userPosts := model.GetUserPosts(svc.DB, user.ID)
	//fmt.Println(userPosts)
	// check add button rendering
	if (username != fmt.Sprintf("%v", sessionUser)) &&
		!(model.CheckFriends(svc.DB, username, fmt.Sprintf("%v", sessionUser))) {
		add = true
	}

	data := map[string]interface{}{
		"user_id":      user.ID,
		"username":     username,
		"name":         user.FirstName,
		"second_name":  user.SecondName,
		"sex":          sex,
		"city":         user.City,
		"interests":    user.Interests,
		"urls":         model.GetFriends(svc.DB, username),
		"add":          add,
		"add_post":     addPost,
		"session_user": sessionUser,
		"posts":        userPosts,
	}
	tmp, err := template.ParseFiles( "web/template/login/user.html")
	if err != nil {
		fmt.Println(err)
	}
	tmp.Execute(resp, data)

}

func AddFriend(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	username := req.Form.Get("username")
	sessionUser := req.Form.Get("user_session")
	_ = model.AddFriend(svc.DB, username, sessionUser)

	user, _  := model.GetUser(svc.DB, username)
	feed := model.GetUserFeed(svc.DB, user.ID)
	data, err := json.Marshal(feed)
	if err != nil {
		fmt.Println("Marshalling feed error")
	}
	ctx := context.Background()
	svc.RDB.Set(ctx, strconv.Itoa(int(user.ID)), data, 0)


	http.Redirect(resp, req, "/", http.StatusSeeOther)
}

func Logout(resp http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "mysession")
	session.Options.MaxAge = -1
	session.Save(req, resp)
	http.Redirect(resp, req, "/", http.StatusSeeOther)
}

func SignUpIndex(resp http.ResponseWriter, req *http.Request) {
	tmp, err := template.ParseFiles("web/template/signup/index.html")
	tmp.Execute(resp, err)
}

func SignUp(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
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
	//fmt.Println(firstName, secondName)

	users := model.NameSearch(svc.DB, firstName[0], secondName[0])

	js, err := json.Marshal(users)
	if err != nil {
		fmt.Println("Users marshalling error")
	}

	resp.Write(js)
}

func Search(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	searchStr := req.Form.Get("search_string")
	userNames := model.SearchAll(svc.DB, searchStr)
	fmt.Println(userNames)
	data := map[string]interface{}{
		"usernames": userNames,
	}
	tmp, _ := template.ParseFiles("web/template/search/results.html")
	tmp.Execute(resp, data)

}
