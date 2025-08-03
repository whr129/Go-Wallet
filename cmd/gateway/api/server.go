package api

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	proxy "github.com/whr129/go-wallet/cmd/gateway/proxy"
	utilLocal "github.com/whr129/go-wallet/cmd/gateway/util"
	"github.com/whr129/go-wallet/pkg/middleware"
	token "github.com/whr129/go-wallet/pkg/token"
	"github.com/whr129/go-wallet/pkg/util"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config      utilLocal.Config
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
	protected.Use(func(ctx *gin.Context) {
		userID, _ := ctx.Get("X-User-ID")
		email, _ := ctx.Get(util.X_EMAIL)
		role, _ := ctx.Get(util.X_ROLE)

		log.Printf("Incoming request - user_id=%v, email=%v, role=%v", userID, email, role)

		ctx.Next() // important to continue to the next handler
	})

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
