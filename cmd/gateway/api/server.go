package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	db "github.com/whr129/go-wallet/cmd/auth-service/db/sqlc"
	proxy "github.com/whr129/go-wallet/cmd/gateway/proxy"
	utilLocal "github.com/whr129/go-wallet/cmd/gateway/util"
	"github.com/whr129/go-wallet/pkg/middleware"
	token "github.com/whr129/go-wallet/pkg/token"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config      utilLocal.Config
	store       db.Store
	tokenMaker  token.Maker
	redisClient *redis.Client
	Router      *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config utilLocal.Config, redisClient *redis.Client) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:      config,
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

	authProxy := proxy.NewReverseProxy(server.config.AuthURL)
	walletProxy := proxy.NewReverseProxy(server.config.WalletURL)

	router.Any("/auth/*path", authProxy)

	protected := router.Group("/").Use(middleware.AuthMiddleware(server.tokenMaker, server.redisClient))
	{
		protected.Any("/wallet/*path", walletProxy)
	}
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
