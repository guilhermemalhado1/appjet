package handlers

import (
	"appjet-decision-manager/app/services"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func InspectAllClustersAllServersHandler(c *gin.Context) {
	config := services.GenerateConfigIfNotExist(c)
	var daemonResponses []map[string]interface{}

	for cIndex := range config.Clusters {
		clusterResponses := make([]map[string]interface{}, 0)

		for sIndex := range config.Clusters[cIndex].Servers {
			serverIP := config.Clusters[cIndex].Servers[sIndex].IP
			daemonResponse, err := services.ForwardInspectToDaemon("http://" + serverIP + ":8080/api/inspect")
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

		clusterResponse := map[string]interface{}{
			config.Clusters[cIndex].Name: clusterResponses,
		}

		daemonResponses = append(daemonResponses, clusterResponse)
	}

	finalResponse := map[string]interface{}{
		"daemon-responses": daemonResponses,
	}

	c.JSON(http.StatusOK, finalResponse)
}

// InspectSpecificClusterAllServersHandler handles inspection of a specific cluster and all servers
func InspectSpecificClusterAllServersHandler(c *gin.Context) {
	config := services.GenerateConfigIfNotExist(c)
	cluster := c.Param("cluster")
	var daemonResponses []map[string]interface{}

	for cIndex := range config.Clusters {
		if config.Clusters[cIndex].Name == cluster {
			clusterResponses := make([]map[string]interface{}, 0)

			for sIndex := range config.Clusters[cIndex].Servers {
				serverIP := config.Clusters[cIndex].Servers[sIndex].IP
				daemonResponse, err := services.ForwardInspectToDaemon("http://" + serverIP + ":8080/api/inspect")
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				var responseBody map[string]interface{}
				err = json.NewDecoder(daemonResponse.Body).Decode(&responseBody)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response body"})
					return
				}
				defer daemonResponse.Body.Close()

				serverResponse := map[string]interface{}{
					config.Clusters[cIndex].Servers[sIndex].Name: responseBody,
				}

				clusterResponses = append(clusterResponses, serverResponse)
			}

			clusterResponse := map[string]interface{}{
				config.Clusters[cIndex].Name: clusterResponses,
			}

			daemonResponses = append(daemonResponses, clusterResponse)
		}
	}

	finalResponse := map[string]interface{}{
		"daemon-responses": daemonResponses,
	}

	c.JSON(http.StatusOK, finalResponse)
}

// InspectSpecificClusterSpecificServerHandler handles inspection of a specific cluster and specific server
func InspectSpecificClusterSpecificServerHandler(c *gin.Context) {
	config := services.GenerateConfigIfNotExist(c)
	cluster := c.Param("cluster")
	server := c.Param("server")
	var daemonResponses []map[string]interface{}

	for cIndex := range config.Clusters {
		if config.Clusters[cIndex].Name == cluster {
			for sIndex := range config.Clusters[cIndex].Servers {
				if config.Clusters[cIndex].Servers[sIndex].Name == server {
					serverIP := config.Clusters[cIndex].Servers[sIndex].IP
					daemonResponse, err := services.ForwardInspectToDaemon("http://" + serverIP + ":8080/api/inspect")

					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					defer daemonResponse.Body.Close()

					var responseBody map[string]interface{}
					err = json.NewDecoder(daemonResponse.Body).Decode(&responseBody)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response body"})
						return
					}

					serverResponse := map[string]interface{}{
						config.Clusters[cIndex].Servers[sIndex].Name: responseBody,
					}

					daemonResponses = append(daemonResponses, serverResponse)
				}
			}
		}
	}

	finalResponse := map[string]interface{}{
		"daemon-responses": daemonResponses,
	}

	c.JSON(http.StatusOK, finalResponse)
}
