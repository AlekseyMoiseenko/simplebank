package api

import (
	"os"
	"testing"
	"time"

	db "github.com/AlekseyMoiseenko/simplebank/db/sqlc"
	"github.com/AlekseyMoiseenko/simplebank/util"
	"github.com/gin-gonic/gin"
)

func newTestServer(store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server := NewServer(config, store)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
