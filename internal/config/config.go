package config

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	DB *sqlx.DB
}

func GetSvc() *Service {
	// TODO: replace with viper
	dbDriver := "mysql"
	dbName := "mydb"//"db"
	dbUser := "root"//"user"
	dbPass := "111"//"password"

	// db, err := sqlx.Open(dbDriver, dbUser+":"+dbPass+"@"+"(db:3306)"+"/"+dbName)
	db, err := sqlx.Open(dbDriver, dbUser+":"+dbPass+"@"+"(0.0.0.0:5506)"+"/"+dbName)
	if err != nil {
		fmt.Println(err)
		fmt.Println("connecting to compose db")
		db, err = sqlx.Open(dbDriver, dbUser+":"+dbPass+"@"+"(db:3306)"+"/"+dbName)
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
	return &Service{
		DB: db,
	}
}
