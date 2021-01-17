package config

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	DB *sqlx.DB
}

func GetDB() *Service {
	// TODO: replace with viper
	dbDriver := "mysql"
	dbName := "db"
	dbUser := "user"
	dbPass := "password"

	db, err := sqlx.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)

	if err != nil {
		panic(err)
	}
	return &Service{
		DB: db,
	}
}
