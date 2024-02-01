package services

import (
	"appjet-server-daemon/app/models"
	"github.com/gin-gonic/gin"
)

func Configure(c *gin.Context) {

	var config models.Configuration

	ShouldBindJson(c, &config)

	GenerateConfigFile(config, c)

	GenerateWaitForItScript(c)

	GenerateDockerfile(config, c)

	GenerateDockerCompose(config, c)

	GenerateInitSqlIfNeeded(config, c)

	c.JSON(200, gin.H{"message": "Configuration saved successfully"})
}
