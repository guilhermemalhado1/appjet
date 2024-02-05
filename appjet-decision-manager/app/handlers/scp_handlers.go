package handlers

import (
	"appjet-decision-manager/app/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SCPCodeAllClustersAllServersHandler(c *gin.Context) {
	config := services.GenerateConfigIfNotExist(c)

	// Slice to store daemon responses
	var daemonResponses []map[string]interface{}

	for cIndex := range config.Clusters {
		clusterResponses := make([]map[string]interface{}, 0)

		for sIndex := range config.Clusters[cIndex].Servers {
			serverIP := config.Clusters[cIndex].Servers[sIndex].IP
			daemonResponse, _ := services.ForwardSCPCodeToDaemon("./code", "http://"+serverIP+":8080/code")

			// Build response structure
			serverResponse := map[string]interface{}{
				"" + config.Clusters[cIndex].Servers[sIndex].Name + "": daemonResponse, // Replace this with the actual response from the daemon
			}

			clusterResponses = append(clusterResponses, serverResponse)
		}

		// Build response structure for the cluster
		clusterResponse := map[string]interface{}{
			"" + config.Clusters[cIndex].Name + "": clusterResponses,
		}

		daemonResponses = append(daemonResponses, clusterResponse)
	}

	// Build the final response structure
	finalResponse := map[string]interface{}{
		"daemon-responses": daemonResponses,
	}

	// Return the final response
	c.JSON(http.StatusOK, finalResponse)
}

func SCPCodeSpecificClusterAllServersHandler(c *gin.Context) {
	config := services.GenerateConfigIfNotExist(c)
	cluster := c.Param("cluster")

	// Slice to store daemon responses
	var daemonResponses []map[string]interface{}

	for cIndex := range config.Clusters {
		if config.Clusters[cIndex].Name == cluster {
			clusterResponses := make([]map[string]interface{}, 0)

			for sIndex := range config.Clusters[cIndex].Servers {
				serverIP := config.Clusters[cIndex].Servers[sIndex].IP
				daemonResponse, _ := services.ForwardSCPCodeToDaemon("./code", "http://"+serverIP+":8080/code")

				// Build response structure
				serverResponse := map[string]interface{}{
					"" + config.Clusters[cIndex].Servers[sIndex].Name + "": daemonResponse, // Replace this with the actual response from the daemon
				}

				clusterResponses = append(clusterResponses, serverResponse)
			}

			// Build response structure for the cluster
			clusterResponse := map[string]interface{}{
				"" + config.Clusters[cIndex].Name + "": clusterResponses,
			}

			daemonResponses = append(daemonResponses, clusterResponse)
		}
	}

	// Build the final response structure
	finalResponse := map[string]interface{}{
		"daemon-responses": daemonResponses,
	}

	// Return the final response
	c.JSON(http.StatusOK, finalResponse)
}

func SCPCodeSpecificClusterSpecificServerHandler(c *gin.Context) {
	config := services.GenerateConfigIfNotExist(c)
	cluster := c.Param("cluster")
	server := c.Param("server")

	// Slice to store daemon responses
	var daemonResponses []map[string]interface{}

	for cIndex := range config.Clusters {
		if config.Clusters[cIndex].Name == cluster {
			for sIndex := range config.Clusters[cIndex].Servers {
				if config.Clusters[cIndex].Servers[sIndex].Name == server {
					serverIP := config.Clusters[cIndex].Servers[sIndex].IP
					daemonResponse, _ := services.ForwardSCPCodeToDaemon("./code", "http://"+serverIP+":8080/code")

					// Build response structure for the server
					serverResponse := map[string]interface{}{
						"" + config.Clusters[cIndex].Servers[sIndex].Name + "": daemonResponse, // Replace this with the actual response from the daemon
					}

					daemonResponses = append(daemonResponses, serverResponse)
				}
			}
		}
	}

	// Build the final response structure
	finalResponse := map[string]interface{}{
		"daemon-responses": daemonResponses,
	}

	// Return the final response
	c.JSON(http.StatusOK, finalResponse)
}
