package dialogServ

import (
	"context"
	dialog "social_network/api/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *dialogServer) GetMessages(ctx context.Context, req *dialog.GetMessagesRequest) (*dialog.GetMessagesResponse, error) {
	var (
		idUser1, idUser2 int64
		message, date    string
		res              []*dialog.Message
	)
	sqlQuery := `SELECT id_user_1, id_user_2, message, dttm_inserted FROM dialogs 
	WHERE (id_user_1=$1 AND id_user_2=$2) OR (id_user_1=$2 AND id_user_2=$1)  ORDER BY dttm_inserted LIMIT $3;`
	rows, err := s.db.Query(sqlQuery, req.IdUser_1, req.IdUser_2, req.Limit)
	if err != nil {
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&idUser1, &idUser2, &message, &date); err != nil {
		}
		res = append(res, &dialog.Message{
			IdUser_1: idUser1,
			IdUser_2: idUser2,
			Message:  message,
			Date:     date,
		})
	}

	return &dialog.GetMessagesResponse{Messages: res}, err

}

func (s *dialogServer) WriteMessage(ctx context.Context, req *dialog.WriteMessageRequest) (*emptypb.Empty, error) {
	sqlQuery := `INSERT INTO dialogs (id_user_1, id_user_2, message) VALUES ($1, $2, $3);`
	_, err := s.db.Query(sqlQuery, req.IdUser_1, req.IdUser_2, req.Message)
	if err != nil {}
	return &empty.Empty{}, nil
}
