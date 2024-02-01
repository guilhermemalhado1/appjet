package services

import (
	"appjet-server-daemon/app/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

func GenerateConfigFile(config models.Configuration, c *gin.Context) {
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

func GenerateWaitForItScript(c *gin.Context) {
	waitForItScriptContent := GenerateWaitForItScriptContent()
	err := CreateFile("wait-for-it.sh", waitForItScriptContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Docker Compose file"})
		return
	}
}

func CreateFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}
