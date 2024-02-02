// main.go
package main

import (
	handlers "appjet-decision-manager/app/handlers"
	services "appjet-decision-manager/app/services"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	_, err := services.CreateDbConnection()
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
		protectedGroup.Use(handlers.AuthMiddlewareHandler)
		{

			//returns config and builds all dependencies - but don't start the process - in all servers in all clusters
			protectedGroup.POST("/configure", handlers.ConfigureAllClustersAllServersHandler)
			//returns config and builds all dependencies - but don't start the process - in all servers in specific cluster
			protectedGroup.POST("/configure/:cluster", handlers.ConfigureSpecificClusterAllServersHandler)
			//returns config and builds all dependencies - but don't start the process - in specific server in specific cluster
			protectedGroup.POST("/configure/:cluster/:server", handlers.ConfigureSpecificClusterSpecificServerHandler)

			/*//start all infrastructure in all servers on the clusters
			apiGroup.GET("/start", handlers.StartHandler)
			//start all infrastructure in all servers on specific cluster
			apiGroup.GET("/start/:cluster", handlers.StartHandler)
			//start all infrastructure in a specific server on specific cluster
			apiGroup.GET("/start/:cluster/:server", handlers.StartHandler)
			//start a specific docker container inside a specific server on specific cluster
			apiGroup.GET("/start/:cluster/:server/:container", handlers.StartContainerHandler)

			//restart all infrastructure in all servers in all clusters
			apiGroup.GET("/restart", handlers.RestartHandler)
			//restart all infrastructure in all servers on specific cluster
			apiGroup.GET("/restart/:cluster", handlers.RestartHandler)
			//restart all infrastructure in a specific server on specific cluster
			apiGroup.GET("/restart/:cluster/:server", handlers.RestartHandler)
			//restart a specific docker container inside a specific server on specific cluster
			apiGroup.GET("/restart/:cluster/:server/:container", handlers.RestartContainerHandler)

			//stop all infrastructure in all servers on the clusters
			apiGroup.GET("/stop", handlers.StopHandler)
			//stop all infrastructure in all servers on specific cluster
			apiGroup.GET("/stop/:cluster", handlers.StopHandler)
			//stop all infrastructure in a specific server on specific cluster
			apiGroup.GET("/stop/:cluster/:server", handlers.StopHandler)
			//stop a specific docker container inside a specific server on specific cluster
			apiGroup.GET("/stop/:cluster/:server/:container", handlers.StopContainerHandler)

			//Check if all containers are alive in all servers in all clusters
			protectedGroup.GET("/check-alive", handlers.CheckAliveAllClustersAllServers)
			//Check if all containers are alive in all servers in specific cluster
			protectedGroup.GET("/check-alive/:cluster", handlers.CheckAliveInClusterAllServers)
			//Check if all containers are alive in specific server in specific cluster
			protectedGroup.GET("/check-alive/:cluster/:server", handlers.CheckAliveInClusterInServer)

			//returns the config.json present in all servers on all clusters
			apiGroup.GET("/inspect", handlers.InspectHandler)
			//returns the config.json present in all servers on a specific clusters
			apiGroup.GET("/inspect/:cluster", handlers.InspectHandler)
			//returns the config.json present in specific server on a specific cluster
			apiGroup.GET("/inspect/:cluster/:server", handlers.InspectHandler)

			//clean all docker images, containers and volumes in all servers in all clusters
			apiGroup.GET("/clean", handlers.CleanHandler)
			//clean all docker images, containers and volumes in all servers in specific clusters
			apiGroup.GET("/clean/:cluster", handlers.CleanHandler)
			//clean all docker images, containers and volumes in specific server in specific clusters
			apiGroup.GET("/clean/:cluster/:server", handlers.CleanHandler)

			//endpoint to load scrips files throught SCP in all servers in all clusters
			apiGroup.POST("/scripts", handlers.SCPHandler)
			//endpoint to load scrips files throught SCP in all servers in specific cluster
			apiGroup.POST("/scripts/:cluster", handlers.SCPHandler)
			//endpoint to load scrips files throught SCP in specific server in specific cluster
			apiGroup.POST("/scripts/:cluster/:server", handlers.SCPHandler)

			//endpoint to load project files throught SCP in all servers in all clusters
			apiGroup.POST("/code", handlers.SCPCodeHandler)
			//endpoint to load project files throught SCP in all servers in specific cluster
			apiGroup.POST("/code/:cluster", handlers.SCPCodeHandler)
			//endpoint to load project files throught SCP in specific server in specific cluster
			apiGroup.POST("/code/:cluster/:server", handlers.SCPCodeHandler)

			//endpoint to run a pre-loaded scp script in all servers in all clusters
			apiGroup.GET("/scp/run/:script", handlers.SCPRunHandler)
			//endpoint to run a pre-loaded scp script in all servers in specific cluster
			apiGroup.GET("/scp/run/:script/:cluster", handlers.SCPRunHandler)
			//endpoint to run a pre-loaded scp script in specific server in specific cluster
			apiGroup.GET("/scp/run/:script/:cluster/:server", handlers.SCPRunHandler)*/
		}
	}

	err = r.Run(":8080")
	if err != nil {
		print(err)
		return
	}
}
