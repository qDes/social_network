package main

import (
	"log"
	"net"
	dialog "social_network/api/proto"
	dialogServ "social_network/internal/app/dialog"

	"google.golang.org/grpc"
)

func main() {
	dbConn := "postgresql://user:pass@0.0.0.0:9118/postgres?sslmode=disable"
	port := "11000"

	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	dialog.RegisterDialogServiceServer(grpcServer, dialogServ.NewServer(dbConn))
	grpcServer.Serve(lis)

}
