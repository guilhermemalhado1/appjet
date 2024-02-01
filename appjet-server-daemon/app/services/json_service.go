package services

import (
	"appjet-server-daemon/app/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func ShouldBindJson(c *gin.Context, config *models.Configuration) {
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
}

func HandleJson(c *gin.Context, config models.Configuration) []byte {

	ShouldBindJson(c, &config)

	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to marshal JSON"})
		return nil
	}

	return configJSON
}
