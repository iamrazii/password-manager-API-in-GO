package main

import (
	"razi/config"
	"razi/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitializeDB()
	router := gin.Default()
	routes.RegisterUserRoutes(router)
	routes.RegisterPasswordRoutes(router)
	router.Run(":1019")
}
