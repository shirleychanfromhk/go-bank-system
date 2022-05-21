package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "simplebank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if valid, ok := binding.Validator.Engine().(*validator.Validate); ok {
		valid.RegisterValidation("currency", validCurrency)
	}

	// Account Endpoints
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.getListAccount)
	router.POST("/accounts", server.createAccount)
	router.PUT("/accounts", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	// Transaction Endpoints
	router.POST("/transactions", server.createTransaction)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
