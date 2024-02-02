package handlers

import (
	"appjet-decision-manager/app/models" // Import your models package
	"appjet-decision-manager/app/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// Assuming you have a User model defined in the "models" package
type User = models.User

func LogoutHandler(c *gin.Context) {
	// Extract the token from the URL parameter
	token := c.Param("token")

	// Check if the token exists in the database
	if !validateToken(token) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Delete the user session with the provided token from the database
	err := services.GetDBConnection().Where("token = ?", token).Delete(&models.UserSession{}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func LoginHandler(c *gin.Context) {
	// Get username and password from the request
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Authenticate the user
	user, err := authenticateUser(username, password)
	if err != nil {
		// Authentication failed
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Assuming the user is authenticated successfully, generate a UUID token
	token := generateToken()

	persistedToken, err := persistToken(user.ID, token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error persisting user session token."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": persistedToken})
}

func authenticateUser(username, password string) (*User, error) {
	var user User
	// Query the database to find the user with the provided username and password
	result := services.GetDBConnection().Where("username = ? AND password = ?", username, password).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func generateToken() string {
	// Generate a new UUID (Version 4)
	uuidToken := uuid.New().String()
	return uuidToken
}

func HelpHandler(c *gin.Context) {
	// Map of commands with descriptions
	commands := map[string]string{
		"./appjet help":             "Show list of available commands",
		"./appjet login <username>": "Authentication",
		"./appjet logout":           "Cancel Authentication",
		"./appjet start":            "Launch config.json and deploy it on target cluster in all servers",
		"./appjet show":             "Should all IP’s inside the configured cluster",
		"./appjet inspect":          "Show deployed application information",
		"./appjet stop":             "Stop Docker process on all IPs of the cluster",
		"./appjet stop <IP>":        "Stop Docker process on desired IP",
		"./appjet restart <IP>":     "Stop Docker process on desired IP and repeat start process only for the desired IP",
		"./appjet restart":          "Stop Docker in all IPs and repeat start process for all IPs",
		"./appjet clean":            "Delete all the running Docker containers in all servers in the cluster",
		"./appjet clean <IP>":       "Delete all the running Docker containers in a specific server in the cluster",
	}

	// Your help logic here
	c.JSON(http.StatusOK, gin.H{"commands": commands})
}

func Endpoint1Handler(c *gin.Context) {
	// Your endpoint1 logic here
	c.JSON(200, gin.H{"message": "Endpoint 1"})
}

func Endpoint2Handler(c *gin.Context) {
	// Your endpoint2 logic here
	c.JSON(200, gin.H{"message": "Endpoint 2"})
}

func validateToken(token string) bool {
	var userSession models.UserSession

	// Query the database to find a user session with the provided token
	result := services.GetDBConnection().Where("token = ?", token).First(&userSession)

	// Check if the token exists in the database
	return result.RowsAffected > 0
}

// Example usage in your middleware or handler
func AuthMiddlewareHandler(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" || !validateToken(token) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	c.Next()
}

func persistToken(userID uint, token string) (string, error) {

	// Create a new UserSession instance
	userSession := models.UserSession{
		UserID: userID,
		Token:  token,
	}

	// Insert the UserSession into the database
	result := services.GetDBConnection().Create(&userSession)

	if result.Error != nil {
		// Error occurred during insert
		return "", fmt.Errorf("error inserting token into the database: %w", result.Error)
	}

	// Token successfully inserted
	return token, nil
}
