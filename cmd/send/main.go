package main

import (
	"context"
	"encoding/json"
	"fmt"
	"social_network/internal/config"
	"social_network/internal/model"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func main() {
	svc := config.GetSvc()
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

	var test model.Feed
	for _, id := range IDS {
		dat := svc.RDB.Get(ctx, strconv.Itoa(int(id)))
		if dat.Err() == redis.Nil {
			fmt.Println("No key")
		}
		data22, err := dat.Bytes()
		if err != nil {
			fmt.Println("UnMarshalling feed error", err)
		}
		json.Unmarshal(data22, &test)
		fmt.Println(test)
	}

	dat := svc.RDB.Get(ctx, "10")
	fmt.Println(dat.Err())
	dat = svc.RDB.Get(ctx, "1")
	fmt.Println(dat.Err())
}
