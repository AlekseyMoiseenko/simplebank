package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/AlekseyMoiseenko/simplebank/api"
	db "github.com/AlekseyMoiseenko/simplebank/db/sqlc"
	"github.com/AlekseyMoiseenko/simplebank/gapi"
	"github.com/AlekseyMoiseenko/simplebank/pb"
	"github.com/AlekseyMoiseenko/simplebank/util"
	_ "github.com/lib/pq"
	_ "go.uber.org/mock/mockgen/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)

	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server := gapi.NewServer(config, store)

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC sever:", err)
	}
}

func runGinServer(config util.Config, store db.Store) {
	server := api.NewServer(config, store)

	err := server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start sever:", err)
	}
}
