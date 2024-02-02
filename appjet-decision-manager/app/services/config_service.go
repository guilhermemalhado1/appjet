package services

import (
	"appjet-decision-manager/app/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
)

func GenerateConfigIfNotExist(c *gin.Context) models.Configuration {
	// Check if config.json file exists
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		// File doesn't exist, generate it with the provided configuration
		var config models.Configuration
		shouldBindJSON(c, &config)
		generateConfigFile(config, c)
		return config
	} else {
		// File exists, read its contents and unmarshal into the config variable
		config, _ := readConfigFile()
		return config
	}
}

func shouldBindJSON(c *gin.Context, config *models.Configuration) {
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
}

func generateConfigFile(config models.Configuration, c *gin.Context) {
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

	c.JSON(200, gin.H{"message": "Config.json generated successfully", "config": config})
}

func readConfigFile() (models.Configuration, error) {
	var config models.Configuration

	// Read the contents of config.json
	configJSON, err := ioutil.ReadFile("config.json")
	if err != nil {
		return config, err
	}

	// Unmarshal the JSON into the config variable
	err = json.Unmarshal(configJSON, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
