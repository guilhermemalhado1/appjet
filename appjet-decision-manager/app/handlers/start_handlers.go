package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)
import services "appjet-decision-manager/app/services"

func StartAllClustersAllServersHandler(c *gin.Context) {
	config := services.GenerateConfigIfNotExist(c)

	// Slice to store daemon responses
	var daemonResponses []map[string]interface{}

	for cIndex := range config.Clusters {
		clusterResponses := make([]map[string]interface{}, 0)

		for sIndex := range config.Clusters[cIndex].Servers {
			serverIP := config.Clusters[cIndex].Servers[sIndex].IP
			daemonResponse, err := services.ForwardStartToDaemon("http://" + serverIP + ":8080/api/start")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

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

			var serverData map[string]interface{}
			err = json.Unmarshal(responseBody, &serverData)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response body"})
				return
			}

			serverResponse := map[string]interface{}{
				config.Clusters[cIndex].Servers[sIndex].Name: serverData,
			}

			clusterResponses = append(clusterResponses, serverResponse)

		}

		// Build response structure for the cluster
		clusterResponse := map[string]interface{}{
			config.Clusters[cIndex].Name: clusterResponses,
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

func StartSpecificClusterAllServersHandler(c *gin.Context) {
	config := services.GenerateConfigIfNotExist(c)
	cluster := c.Param("cluster")

	// Slice to store daemon responses
	var daemonResponses []map[string]interface{}

	for cIndex := range config.Clusters {
		if config.Clusters[cIndex].Name == cluster {
			clusterResponses := make([]map[string]interface{}, 0)

			for sIndex := range config.Clusters[cIndex].Servers {
				serverIP := config.Clusters[cIndex].Servers[sIndex].IP
				daemonResponse, err := services.ForwardStartToDaemon("http://" + serverIP + ":8080/api/start")
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

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

				var serverData map[string]interface{}
				err = json.Unmarshal(responseBody, &serverData)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response body"})
					return
				}

				serverResponse := map[string]interface{}{
					config.Clusters[cIndex].Servers[sIndex].Name: serverData,
				}

				clusterResponses = append(clusterResponses, serverResponse)
			}

			// Build response structure for the cluster
			clusterResponse := map[string]interface{}{
				config.Clusters[cIndex].Name: clusterResponses,
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

func StartSpecificClusterSpecificServerHandler(c *gin.Context) {
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
					daemonResponse, err := services.ForwardStartToDaemon("http://" + serverIP + ":8080/api/start")
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}

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

					var serverData map[string]interface{}
					err = json.Unmarshal(responseBody, &serverData)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response body"})
						return
					}

					serverResponse := map[string]interface{}{
						config.Clusters[cIndex].Servers[sIndex].Name: serverData,
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

func StartContainerSpecificClusterSpecificServerHandler(c *gin.Context) {
	config := services.GenerateConfigIfNotExist(c)
	cluster := c.Param("cluster")
	server := c.Param("server")
	container := c.Param("container")

	// Slice to store daemon responses
	var daemonResponses []map[string]interface{}

	for cIndex := range config.Clusters {
		if config.Clusters[cIndex].Name == cluster {
			for sIndex := range config.Clusters[cIndex].Servers {
				if config.Clusters[cIndex].Servers[sIndex].Name == server {
					serverIP := config.Clusters[cIndex].Servers[sIndex].IP
					daemonResponse, _ := services.ForwardStartToDaemon("http://" + serverIP + ":8080/api/start/" + container)

					// Build response structure for the server
					serverResponse := map[string]interface{}{
						config.Clusters[cIndex].Servers[sIndex].Name: daemonResponse, // Storing the actual response from the daemon
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
