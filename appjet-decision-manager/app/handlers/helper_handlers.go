package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HelpHandler(c *gin.Context) {
	// Map to store endpoint descriptions
	endpointDescriptions := map[string]map[string]string{
		"no-auth-commands": {
			"./appjet login": "For user authentication.",
			"./appjet help":  "To see available appjet commands.",
		},
		"needed-auth-commands": {
			"./appjet logout":                       "To cancel user authentication.",
			"./appjet check-alive":                  "Check if all containers are alive in all servers in all clusters",
			"./appjet check-alive/:cluster":         "Check if all containers are alive in all servers in specific cluster",
			"./appjet check-alive/:cluster/:server": "Check if all containers are alive in specific server in specific cluster",

			"./appjet configure":                  "Load config and install all Docker containers (without starting them)",
			"./appjet configure/:cluster":         "Load config and install all Docker containers in a specific cluster (without starting them)",
			"./appjet configure/:cluster/:server": "Load config and install all Docker containers in a specific server in a specific cluster (without starting them)",

			"./appjet inspect":                  "Return the config.json present in all servers on the cluster",
			"./appjet inspect/:cluster":         "Return the config.json present in all servers on a specific cluster",
			"./appjet inspect/:cluster/:server": "Return the config.json present in a specific server on a specific cluster",

			"./appjet start":                             "Start all infrastructure in all servers on all clusters",
			"./appjet start/:cluster":                    "Start all infrastructure in all servers on a specific cluster",
			"./appjet start/:cluster/:server":            "Start all infrastructure in a specific server on a specific cluster",
			"./appjet start/:cluster/:server/:container": "Start a specific Docker container inside a specific server on a specific cluster",

			"./appjet restart":                             "Restart all infrastructure in all servers in all clusters",
			"./appjet restart/:cluster":                    "Restart all infrastructure in all servers on a specific cluster",
			"./appjet restart/:cluster/:server":            "Restart all infrastructure in a specific server on a specific cluster",
			"./appjet restart/:cluster/:server/:container": "Restart a specific Docker container inside a specific server on a specific cluster",

			"./appjet stop/:cluster":                    "Stop all infrastructure in all servers on a specific cluster",
			"./appjet stop/:cluster/:server":            "Stop all infrastructure in a specific server on a specific cluster",
			"./appjet stop":                             "Stop all infrastructure in all servers on all clusters",
			"./appjet stop/:cluster/:server/:container": "Stop a specific Docker container inside a specific server on a specific cluster",

			"./appjet clean":                  "Clean all Docker images, containers, and volumes in all servers in all clusters",
			"./appjet clean/:cluster":         "Clean all Docker images, containers, and volumes in all servers in a specific cluster",
			"./appjet clean/:cluster/:server": "Clean all Docker images, containers, and volumes in a specific server in a specific cluster",

			"./appjet scripts":                  "Load scripts files through SCP in all servers in all clusters",
			"./appjet scripts/:cluster":         "Load scripts files through SCP in all servers in a specific cluster",
			"./appjet scripts/:cluster/:server": "Load scripts files through SCP in a specific server in a specific cluster",

			"./appjet code":                  "Load project files through SCP in all servers in all clusters",
			"./appjet code/:cluster":         "Load project files through SCP in all servers in a specific cluster",
			"./appjet code/:cluster/:server": "Load project files through SCP in a specific server in a specific cluster",

			"./appjet scp/run/:script":                  "Run a pre-loaded SCP script in all servers in all clusters",
			"./appjet scp/run/:script/:cluster":         "Run a pre-loaded SCP script in all servers in a specific cluster",
			"./appjet scp/run/:script/:cluster/:server": "Run a pre-loaded SCP script in a specific server in a specific cluster",
		},
	}

	c.JSON(http.StatusOK, gin.H{"endpoints": endpointDescriptions})
}
