package gapi

import (
	db "github.com/AlekseyMoiseenko/simplebank/db/sqlc"
	"github.com/AlekseyMoiseenko/simplebank/pb"
	"github.com/AlekseyMoiseenko/simplebank/token"
	"github.com/AlekseyMoiseenko/simplebank/util"
)

// Server serves gRPC requests for our banking service
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) *Server {
	tokenMaker := token.NewPasetoMaker(config.TokenSymmetricKey)
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server
}
