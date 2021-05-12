package main

import (
	"log"
	"net"
	dialog "social_network/api/proto"
	dialogServ "social_network/internal/app/dialog"

	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"
)

func main() {
	// dbConn := "postgresql://user:pass@0.0.0.0:9118/postgres?sslmode=disable"
	dbConn := "postgresql://user:pass@db-dialog:5432/postgres?sslmode=disable"

	port := "11000"

	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_opentracing.UnaryServerInterceptor()))

	dialog.RegisterDialogServiceServer(grpcServer, dialogServ.NewServer(dbConn))
	grpcServer.Serve(lis)

}
