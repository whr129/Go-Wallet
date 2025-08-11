package gapi

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/whr129/go-wallet/cmd/wallet-service/db/sqlc"
	"github.com/whr129/go-wallet/internal/pb"
	"github.com/whr129/go-wallet/pkg/token"
	util "github.com/whr129/go-wallet/pkg/util"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	pb.UnimplementedWalletServiceServer
	config     util.Config
	tokenMaker token.Maker
	store      db.Store
	Router     *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}

// Start runs the HTTP server on a specific address.
// func (server *Server) Start(address string) error {
// 	return server.Router.Run(address)
// }

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
