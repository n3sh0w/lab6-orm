package routes

import (
	"go-postgres-orm/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		// User routes
		api.GET("/users", controllers.GetUsers)
		api.POST("/users", controllers.CreateUser)
		api.GET("/users/:id", controllers.GetUser)
		api.PUT("/users/:id", controllers.UpdateUser)
		api.DELETE("/users/:id", controllers.DeleteUser)

		// Product routes
		api.GET("/products", controllers.GetProducts)
		api.POST("/products", controllers.CreateProduct)
		api.PUT("/products/:id", controllers.UpdateProduct)
		api.DELETE("/products/:id", controllers.DeleteProduct)

		// Order routes
		api.GET("/orders", controllers.GetOrders)
		api.POST("/orders", controllers.CreateOrder)
		api.PUT("/orders/:id", controllers.UpdateOrder)
		api.DELETE("/orders/:id", controllers.DeleteOrder)
	}

	return r
}
