package config

import (
	"fmt"
	dialog "social_network/api/proto"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/streadway/amqp"
	"github.com/tarantool/go-tarantool"
	"google.golang.org/grpc"
)

type Service struct {
	DB           *sqlx.DB
	RDB          *redis.Client
	Feed         *amqp.Channel
	Q            amqp.Queue
	Tarantool    *tarantool.Connection
	DialogClient dialog.DialogServiceClient
}

func GetSvc() *Service {
	// TODO: replace with viper

	// connection to mysql
	dbDriver := "mysql"
	dbName := "db"       //"mydb"//
	dbUser := "user"     //"root"//
	dbPass := "password" //"secret"

	db, err := sqlx.Open(dbDriver, dbUser+":"+dbPass+"@"+"(db-master:3306)"+"/"+dbName+"?parseTime=true")
	//db, err := sqlx.Open(dbDriver, dbUser+":"+dbPass+"@"+"(0.0.0.0:3306)"+"/"+dbName+"?parseTime=true")
	if err != nil {
		fmt.Println(err)
		fmt.Println("connecting to compose db")
		db, err = sqlx.Open(dbDriver, dbUser+":"+dbPass+"@"+"(db:3306)"+"/"+dbName+"?parseTime=true")
		if err != nil {
			panic(err)
		}
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	db.SetMaxOpenConns(2500)
	db.SetMaxIdleConns(2500)
	db.SetConnMaxLifetime(time.Duration(time.Duration.Seconds(1)))

	// connection to rabbitmq
	// conn, err := amqp.Dial("amqp://rabbit:rabbit@0.0.0.0:5672/")
	conn, err := amqp.Dial("amqp://rabbit:rabbit@rabbitmq:5672/")
	if err != nil {
		fmt.Println(err)
	}
	ch, err := conn.Channel()

	q, err := ch.QueueDeclare(
		"feed", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)

	//connection to redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		//DB:       6379,
	})

	//tarantool, err := tarantool.Connect("127.0.0.1:3301", tarantool.Opts{

	tarantool, err := tarantool.Connect("tarantool:3301", tarantool.Opts{
		User: "admin",
		Pass: "admin",
	})

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithInsecure())

	// cn, err := grpc.Dial("0.0.0.0:11000", opts...)
	cn, err := grpc.Dial("dialog:11000", opts...)
	if err != nil {
		panic(err)
	}
	dialogClient := dialog.NewDialogServiceClient(cn)

	return &Service{
		DB:           db,
		RDB:          rdb,
		Feed:         ch,
		Q:            q,
		Tarantool:    tarantool,
		DialogClient: dialogClient,
	}
}
