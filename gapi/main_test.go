package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/AlekseyMoiseenko/simplebank/db/sqlc"
	"github.com/AlekseyMoiseenko/simplebank/token"
	"github.com/AlekseyMoiseenko/simplebank/util"
	"github.com/AlekseyMoiseenko/simplebank/worker"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func newTestServer(store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := &util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server := NewServer(config, store, taskDistributor)

	return server
}

func newContextWithBearerToken(t *testing.T, tokenMaker token.Maker, username string, role string, duration time.Duration) context.Context {
	accessToken, _, err := tokenMaker.CreateToken(username, role, duration)
	require.NoError(t, err)

	bearerToken := fmt.Sprintf("%s %s", authorizationBearer, accessToken)
	md := metadata.MD{
		authorizationHeader: []string{
			bearerToken,
		},
	}

	return metadata.NewIncomingContext(context.Background(), md)
}
