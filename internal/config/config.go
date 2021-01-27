package config

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	DB *sqlx.DB
}

func GetSvc() *Service {
	// TODO: replace with viper
	dbDriver := "mysql"
	dbName := "db"
	dbUser := "user"
	dbPass := "password"

	// db, err := sqlx.Open(dbDriver, dbUser+":"+dbPass+"@"+"(db:3306)"+"/"+dbName)
	db, err := sqlx.Open(dbDriver, dbUser+":"+dbPass+"@"+"(0.0.0.0:3306)"+"/"+dbName)
	if err != nil {
		fmt.Println(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}
	return &Service{
		DB: db,
	}
}
