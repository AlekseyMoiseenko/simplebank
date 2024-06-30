package gapi

import (
	db "github.com/AlekseyMoiseenko/simplebank/db/sqlc"
	"github.com/AlekseyMoiseenko/simplebank/pb"
	"github.com/AlekseyMoiseenko/simplebank/token"
	"github.com/AlekseyMoiseenko/simplebank/util"
	"github.com/AlekseyMoiseenko/simplebank/worker"
)

// Server serves gRPC requests for our banking service
type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	tokenMaker := token.NewPasetoMaker(config.TokenSymmetricKey)
	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server
}
