// main.go
package main

import (
	commandhandler "appjet-server-daemon/app/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	apiGroup := r.Group("/api")
	{
		//check if alive containers endpoint
		apiGroup.GET("/check-alive", commandhandler.CheckAlive)

		//load config and install all docker containers - but don't start them
		apiGroup.POST("/configure", commandhandler.ConfigureHandler)
		//returns the config.json present in the current daemon
		apiGroup.GET("/inspect", commandhandler.InspectHandler)

		//start all configured docker containers
		apiGroup.GET("/start", commandhandler.StartHandler)
		//restart all docker containers
		apiGroup.GET("/restart", commandhandler.RestartHandler)
		//stop all docker containers
		apiGroup.GET("/stop", commandhandler.StopHandler)

		//start a specific docker container
		apiGroup.GET("/start/:container", commandhandler.StartContainerHandler)
		//restart a specific docker container
		apiGroup.GET("/restart/:container", commandhandler.RestartContainerHandler)
		//stop a specific container
		apiGroup.GET("/stop/:container", commandhandler.StopContainerHandler)

		//clean all docker images, containers and volumes in this server (docker system prune -a)
		apiGroup.GET("/clean", commandhandler.CleanHandler)

		//endpoint to run a pre-loaded scp script
		apiGroup.GET("/scp/run/:script", commandhandler.SCPRunHandler)

		//endpoint to load scrips files throught SCP
		apiGroup.POST("/scripts", commandhandler.SCPHandler)

		//endpoint to load project files throught SCP
		apiGroup.POST("/code", commandhandler.SCPCodeHandler)
	}

	err := r.Run(":8080")
	if err != nil {
		print(err)
		return
	}
}
