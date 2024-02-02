// main.go
package main

import (
	handlers "appjet-decision-manager/app/handlers"
	services "appjet-decision-manager/app/services"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
)

func main() {
	r := gin.Default()

	portStr := os.Getenv("port")
	if portStr == "" {
		portStr = "8080"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("Invalid port number:", err)
	}

	_, err = services.CreateDbConnection()
	if err != nil {
		return
	}

	apiGroup := r.Group("/appjet")
	{
		// open endpoints
		apiGroup.POST("/login", handlers.LoginHandler)
		apiGroup.GET("/logout/:token", handlers.LogoutHandler)

		apiGroup.GET("/help", handlers.HelpHandler)

		// protected endpoints
		protectedGroup := apiGroup.Group("/")
		protectedGroup.Use(handlers.AuthMiddlewareHandler) // Add middleware for token authentication
		{
			protectedGroup.GET("/endpoint1", handlers.Endpoint1Handler)
			protectedGroup.GET("/endpoint2", handlers.Endpoint2Handler)
		}
	}

	err = r.Run(":" + strconv.Itoa(port))
	if err != nil {
		print(err)
		return
	}
}
