// main.go
package main

import (
	handlers "appjet-decision-manager/app/handlers"
	services "appjet-decision-manager/app/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"Authorization"} // Allow only Authorization header
	config.AllowOrigins = []string{"*"}             // Allow all origins
	r.Use(cors.New(config))

	_, err := services.CreateDbConnection()
	if err != nil {
		return
	}

	apiGroup := r.Group("/appjet")
	{
		// open endpoints
		apiGroup.POST("/login", handlers.LoginHandler)         //OK
		apiGroup.GET("/logout/:token", handlers.LogoutHandler) //OK

		apiGroup.GET("/help", handlers.HelpHandler) //OK

		// protected endpoints
		protectedGroup := apiGroup.Group("/")
		protectedGroup.Use(handlers.AuthMiddlewareHandler)
		{

			//returns config and builds all dependencies - but don't start the process - in all servers in all clusters
			protectedGroup.POST("/configure", handlers.ConfigureAllClustersAllServersHandler) //OK
			//returns config and builds all dependencies - but don't start the process - in all servers in specific cluster
			protectedGroup.POST("/configure/:cluster", handlers.ConfigureSpecificClusterAllServersHandler) //OK
			//returns config and builds all dependencies - but don't start the process - in specific server in specific cluster
			protectedGroup.POST("/configure/:cluster/:server", handlers.ConfigureSpecificClusterSpecificServerHandler) //OK

			//start all infrastructure in all servers on the clusters
			apiGroup.GET("/start", handlers.StartAllClustersAllServersHandler) //OK
			//start all infrastructure in all servers on specific cluster
			apiGroup.GET("/start/:cluster", handlers.StartSpecificClusterAllServersHandler) //OK
			//start all infrastructure in a specific server on specific cluster
			apiGroup.GET("/start/:cluster/:server", handlers.StartSpecificClusterSpecificServerHandler) //OK
			//start a specific docker container inside a specific server on specific cluster
			apiGroup.GET("/start/:cluster/:server/:container", handlers.StartContainerSpecificClusterSpecificServerHandler)

			//restart all infrastructure in all servers in all clusters
			apiGroup.GET("/restart", handlers.RestartAllClustersAllServersHandler)
			//restart all infrastructure in all servers on specific cluster
			apiGroup.GET("/restart/:cluster", handlers.RestartSpecificClusterAllServersHandler)
			//restart all infrastructure in a specific server on specific cluster
			apiGroup.GET("/restart/:cluster/:server", handlers.RestartSpecificClusterSpecificServerHandler)
			//restart a specific docker container inside a specific server on specific cluster
			apiGroup.GET("/restart/:cluster/:server/:container", handlers.RestartContainerSpecificClusterSpecificServerContainerHandler)

			//stop all infrastructure in all servers on the clusters
			apiGroup.GET("/stop", handlers.StopAllClustersAllServersHandler)
			//stop all infrastructure in all servers on specific cluster
			apiGroup.GET("/stop/:cluster", handlers.StopSpecificClusterAllServersHandler)
			//stop all infrastructure in a specific server on specific cluster
			apiGroup.GET("/stop/:cluster/:server", handlers.StopSpecificClusterSpecificServerHandler)
			//stop a specific docker container inside a specific server on specific cluster
			apiGroup.GET("/stop/:cluster/:server/:container", handlers.StopContainerSpecificClusterSpecificServerContainerHandler)

			//Check if all containers are alive in all servers in all clusters
			protectedGroup.GET("/check-alive", handlers.CheckAliveAllClustersAllServersHandler)
			//Check if all containers are alive in all servers in specific cluster
			protectedGroup.GET("/check-alive/:cluster", handlers.CheckAliveSpecificClusterAllServersHandler)
			//Check if all containers are alive in specific server in specific cluster
			protectedGroup.GET("/check-alive/:cluster/:server", handlers.CheckAliveSpecificClusterSpecificServerHandler)

			//returns the config.json present in all servers on all clusters
			apiGroup.GET("/inspect", handlers.InspectAllClustersAllServersHandler)
			//returns the config.json present in all servers on a specific clusters
			apiGroup.GET("/inspect/:cluster", handlers.InspectSpecificClusterAllServersHandler)
			//returns the config.json present in specific server on a specific cluster
			apiGroup.GET("/inspect/:cluster/:server", handlers.InspectSpecificClusterSpecificServerHandler)

			//clean all docker images, containers and volumes in all servers in all clusters
			apiGroup.GET("/clean", handlers.CleanAllClustersAllServersHandler)
			//clean all docker images, containers and volumes in all servers in specific clusters
			apiGroup.GET("/clean/:cluster", handlers.CleanSpecificClusterAllServersHandler)
			//clean all docker images, containers and volumes in specific server in specific clusters
			apiGroup.GET("/clean/:cluster/:server", handlers.CleanSpecificClusterSpecificServerHandler)

			//endpoint to load scrips files throught SCP in all servers in all clusters
			apiGroup.POST("/scripts", handlers.SCPAllClustersAllServersHandler)
			//endpoint to load scrips files throught SCP in all servers in specific cluster
			apiGroup.POST("/scripts/:cluster", handlers.SCPSpecificClusterAllServersHandler)
			//endpoint to load scrips files throught SCP in specific server in specific cluster
			apiGroup.POST("/scripts/:cluster/:server", handlers.SCPSpecificClusterSpecificServerHandler)

			//endpoint to load project files throught SCP in all servers in all clusters
			apiGroup.POST("/code", handlers.SCPCodeAllClustersAllServersHandler)
			//endpoint to load project files throught SCP in all servers in specific cluster
			apiGroup.POST("/code/:cluster", handlers.SCPCodeSpecificClusterAllServersHandler)
			//endpoint to load project files throught SCP in specific server in specific cluster
			apiGroup.POST("/code/:cluster/:server", handlers.SCPCodeSpecificClusterSpecificServerHandler)

			//endpoint to run a pre-loaded scp script in all servers in all clusters
			apiGroup.GET("/scp/run/:script", handlers.SCPRunAllClustersAllServersHandler)
			//endpoint to run a pre-loaded scp script in all servers in specific cluster
			apiGroup.GET("/scp/run/:script/:cluster", handlers.SCPRunSpecificClusterAllServersHandler)
			//endpoint to run a pre-loaded scp script in specific server in specific cluster
			apiGroup.GET("/scp/run/:script/:cluster/:server", handlers.SCPRunSpecificClusterSpecificServerHandler)
		}
	}

	err = r.Run(":8080")
	if err != nil {
		print(err)
		return
	}
}
