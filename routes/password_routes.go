package routes

import (
	"razi/controllers"
	"razi/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterPasswordRoutes(router *gin.Engine) {
	password := router.Group("/password")
	password.Use(middleware.JWTAuthMiddleware()) // attaching middleware to route group for protection

	{
		password.POST("/create", controllers.StorePassword)

		password.GET("/all", controllers.GetAllPasswords)

		password.PUT("/update/:id", controllers.UpdatePassword)

		password.DELETE("/delete/:id", controllers.DeletePassword)
	}
}
