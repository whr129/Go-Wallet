package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/whr129/go-wallet/cmd/transaction-service/db/sqlc"
	"github.com/whr129/go-wallet/cmd/transaction-service/middleware"
	"github.com/whr129/go-wallet/pkg/util"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config util.Config
	store  db.Store
	Router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
	}

	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidation("currency", validCurrency)
	// }

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.Use(middleware.AuthMiddleware)

	router.POST("/transfers", server.createTransfer)

	server.Router = router
}

// Start runs the HTTP server on a specific address.
// func (server *Server) Start(address string) error {
// 	return server.Router.Run(address)
// }

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
