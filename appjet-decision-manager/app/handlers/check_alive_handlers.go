package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)
import services "appjet-decision-manager/app/services"

func CheckAliveAllClustersAllServersHandler(c *gin.Context) {
	config := services.GenerateConfigIfNotExist(c)

	// Slice to store daemon responses
	var daemonResponses []map[string]interface{}

	for cIndex := range config.Clusters {
		clusterResponses := make([]map[string]interface{}, 0)

		for sIndex := range config.Clusters[cIndex].Servers {
			serverIP := config.Clusters[cIndex].Servers[sIndex].IP
			daemonResponse, err := services.ForwardCheckAliveToDaemon("http://" + serverIP + ":8080/api/check-alive")
			if err != nil {
				// Handle error if needed
				fmt.Println("Error:", err)
				continue
			}

			// Read the response body
			responseBody, err := ioutil.ReadAll(daemonResponse.Body)
			if err != nil {
				// Handle error if needed
				fmt.Println("Error reading response body:", err)
				continue
			}

			// Parse the response body as JSON
			var jsonResponse map[string]interface{}
			if err := json.Unmarshal(responseBody, &jsonResponse); err != nil {
				// Handle error if needed
				fmt.Println("Error parsing JSON:", err)
				continue
			}

			// Build response structure
			serverResponse := map[string]interface{}{
				"server_name": config.Clusters[cIndex].Servers[sIndex].Name,
				"response":    jsonResponse,
			}

			clusterResponses = append(clusterResponses, serverResponse)
		}

		// Build response structure for the cluster
		clusterResponse := map[string]interface{}{
			"cluster_name": config.Clusters[cIndex].Name,
			"servers":      clusterResponses,
		}

		daemonResponses = append(daemonResponses, clusterResponse)
	}

	// Build the final response structure
	finalResponse := map[string]interface{}{
		"daemon_responses": daemonResponses,
	}

	// Return the final response
	c.JSON(http.StatusOK, finalResponse)
}

func CheckAliveSpecificClusterAllServersHandler(c *gin.Context) {
	config := services.GenerateConfigIfNotExist(c)
	cluster := c.Param("cluster")

	// Slice to store daemon responses
	var daemonResponses []map[string]interface{}

	for cIndex := range config.Clusters {
		if config.Clusters[cIndex].Name == cluster {
			clusterResponses := make([]map[string]interface{}, 0)

			for sIndex := range config.Clusters[cIndex].Servers {
				serverIP := config.Clusters[cIndex].Servers[sIndex].IP
				daemonResponse, err := services.ForwardCheckAliveToDaemon("http://" + serverIP + ":8080/api/check-alive")
				if err != nil {
					// Handle error if needed
					fmt.Println("Error:", err)
					continue
				}
				defer daemonResponse.Body.Close()

				// Read the response body
				responseBody, err := ioutil.ReadAll(daemonResponse.Body)
				if err != nil {
					// Handle error if needed
					fmt.Println("Error reading response body:", err)
					continue
				}

				// Build response structure for the server
				serverResponse := map[string]interface{}{
					"server_name": config.Clusters[cIndex].Servers[sIndex].Name,
					"response":    string(responseBody),
				}

				clusterResponses = append(clusterResponses, serverResponse)
			}

			// Build response structure for the cluster
			clusterResponse := map[string]interface{}{
				"cluster_name": config.Clusters[cIndex].Name,
				"servers":      clusterResponses,
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

func CheckAliveSpecificClusterSpecificServerHandler(c *gin.Context) {
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
					daemonResponse, err := services.ForwardCheckAliveToDaemon("http://" + serverIP + ":8080/api/check-alive")
					if err != nil {
						// Handle error if needed
						fmt.Println("Error:", err)
						continue
					}
					defer daemonResponse.Body.Close()

					// Read the response body
					responseBody, err := ioutil.ReadAll(daemonResponse.Body)
					if err != nil {
						// Handle error if needed
						fmt.Println("Error reading response body:", err)
						continue
					}

					// Build response structure for the server
					serverResponse := map[string]interface{}{
						"server_name": config.Clusters[cIndex].Servers[sIndex].Name,
						"response":    string(responseBody),
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
