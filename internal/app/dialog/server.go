package dialogServ

import (
	dialog "social_network/api/proto"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type dialogServer struct {
	dialog.UnimplementedDialogServiceServer
	db *sqlx.DB
}

func NewServer(dbConn string) *dialogServer {
	db := sqlx.MustOpen("postgres", dbConn)
	s := &dialogServer{db: db}
	return s
}
