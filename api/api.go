package api

import (
	"accounting/api/handler"
	"accounting/ent"
	"accounting/service"

	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the Gin router and returns an instance of the router
func SetupRouter(client *ent.Client) *gin.Engine {
	r := gin.Default()

	// Initialize services and handlers
	userService := service.NewUserService(client)
	userHandler := handler.NewUserHandler(userService)

	transactionService := service.NewTransactionService(client)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// API endpoints group
	api := r.Group("/api")
	{
		// Users endpoints
		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
		}

		// Transactions endpoints
		transactions := api.Group("/transactions")
		{
			transactions.POST("", transactionHandler.CreateTransaction)
		}
	}

	return r
}
