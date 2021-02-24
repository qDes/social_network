package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"social_network/internal/model"
	"strconv"

	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)

func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}


func QReader(){
	var post model.Post

	msgs, err := svc.Feed.Consume(
		svc.Q.Name,
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args)
	)
	if err != nil {
		fmt.Println("Qreader error", err)
	}
	for d := range msgs {
		err := json.Unmarshal(d.Body, &post)
		if err != nil {
			fmt.Println("Unmarshalling error", err)
		}
		fmt.Println("RECV:", post.FriendID)
		FeedUpdater(post)

	}
}

func FeedUpdater(post model.Post) {
	var feed model.Feed
	ctx := context.Background()
	// get feed from redis
	data := svc.RDB.Get(ctx, strconv.Itoa(int(post.FriendID)))
	if data.Err() == redis.Nil {
		fmt.Println("No key")
		feed = model.GetUserFeed(svc.DB, post.FriendID)
	} else {
		redisBytes, err := data.Bytes()
		if err != nil {
			fmt.Println("redis bytes error")
		}
		// unmarshal feed
		err = json.Unmarshal(redisBytes, &feed)
		if err != nil {
			fmt.Println("UnMarshalling feed error")
		}
		feed.Posts = append([]model.Post{post}, feed.Posts...)
		if len(feed.Posts) > 1000 {
			feed.Posts = feed.Posts[:1000]
		}
	}

	js, err := json.Marshal(feed)
	if err != nil {
		fmt.Println("marshalling error")
	}
	svc.RDB.Set(ctx, strconv.Itoa(int(post.FriendID)), js,0)

}

func InitFeedCache() {
	IDS := model.GetUsersIDs(svc.DB)
	ctx := context.Background()
	for _, id := range IDS {
		feed := model.GetUserFeed(svc.DB, id)
		data, err := json.Marshal(feed)
		if err != nil {
			fmt.Println("Marshalling feed error")
		}

		svc.RDB.Set(ctx, strconv.Itoa(int(id)), data, 0)
	}
}