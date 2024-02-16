package services

import (
	"appjet-decision-manager/app/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
)

func GenerateConfigIfNotExist(c *gin.Context) *models.Configuration {
	// Check if config.json file exists
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		// File doesn't exist, generate it with the provided configuration
		var config models.Configuration
		shouldBindJSON(c, &config)
		generateConfigFile(&config, c)
		return &config
	} else {
		// File exists, read its contents and unmarshal into the config variable
		config, err := readConfigFile()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to load config.json"})
		}
		return config
	}
}

func shouldBindJSON(c *gin.Context, config *models.Configuration) {
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
}

func generateConfigFile(config *models.Configuration, c *gin.Context) {
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to marshal JSON"})
		return
	}

	err = ioutil.WriteFile("config.json", configJSON, 0644)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to write JSON to file"})
		return
	}

}

func readConfigFile() (*models.Configuration, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the content of config.json
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Parse the content into a Config struct
	var config models.Configuration
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
