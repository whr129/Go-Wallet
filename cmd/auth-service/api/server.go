package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	db "github.com/whr129/go-wallet/cmd/auth-service/db/sqlc"
	token "github.com/whr129/go-wallet/pkg/token"
	util "github.com/whr129/go-wallet/pkg/util"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config      util.Config
	store       db.Store
	tokenMaker  token.Maker
	redisClient *redis.Client
	Router      *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store, redisClient *redis.Client) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:      config,
		store:       store,
		redisClient: redisClient,
		tokenMaker:  tokenMaker,
	}

	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidation("currency", validCurrency)
	// }

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/create", server.createUser)
	router.POST("/login", server.loginUser)
	// router.POST("/tokens/renew_access", server.renewAccessToken)

	server.Router = router
}

// Start runs the HTTP server on a specific address.
// func (server *Server) Start(address string) error {
// 	return server.Router.Run(address)
// }

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
