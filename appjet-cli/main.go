// main.go
package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/login", nil)
	}

	err := r.Run(":8080")
	if err != nil {
		print(err)
		return
	}
}
