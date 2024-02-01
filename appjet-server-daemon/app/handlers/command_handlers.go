package handlers

import (
	"appjet-server-daemon/app/models"
	configurationService "appjet-server-daemon/app/services"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func ConfigureHandler(c *gin.Context) {
	configurationService.Configure(c)
}

func InspectHandler(c *gin.Context) {
	config, err := loadConfig("config.json")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to load config.json"})
		return
	}

	containerStates, err := configurationService.GetDockerContainersState()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get Docker container states"})
		return
	}

	response := gin.H{
		"docker": containerStates,
		"config": config,
	}

	c.JSON(200, response)
}

func loadConfig(filename string) (*models.Configuration, error) {
	file, err := os.Open(filename)
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

func StartHandler(c *gin.Context) {
	// Run "docker compose up -d" command in the current directory
	cmd := exec.Command("docker", "compose", "up", "--build")
	cmd.Dir = "." // Set the command's working directory to the current directory

	// Start the command asynchronously
	err := cmd.Start()
	if err != nil {
		// If there's an error starting the command, respond with an error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to start Docker Compose: %s", err.Error())})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Docker Compose started successfully"})

}

func RestartHandler(c *gin.Context) {
	containerNames := []string{"app", "database"}
	for _, name := range containerNames {
		err := restartContainer(name)
		if err != nil {
			log.Printf("Error restarting container %s: %s", name, err)
			c.JSON(500, gin.H{"error": "Failed to restart containers"})
			return
		}
	}
	c.JSON(200, gin.H{"message": "Containers restarted successfully"})
}

func StopHandler(c *gin.Context) {
	containerNames := []string{"app", "database"}
	for _, name := range containerNames {
		err := stopContainer(name)
		if err != nil {
			log.Printf("Error stopping container %s: %s", name, err)
			c.JSON(500, gin.H{"error": "Failed to stop containers"})
			return
		}
	}
	c.JSON(200, gin.H{"message": "Containers stopped successfully"})
}

func restartContainer(containerName string) error {
	cmd := exec.Command("docker", "container", "restart", containerName)
	err := cmd.Run()
	return err
}

func stopContainer(containerName string) error {
	cmd := exec.Command("docker", "container", "stop", containerName)
	err := cmd.Run()
	return err
}

func startContainer(containerName string) error {
	cmd := exec.Command("docker", "container", "start", containerName)
	err := cmd.Run()
	return err
}

func StartContainerHandler(c *gin.Context) {
	// Retrieve the container parameter from the URL path
	container := c.Param("container")
	err := startContainer(container)
	if err != nil {
		log.Printf("Error stopping container %s: %s", container, err)
		c.JSON(500, gin.H{"error": "Failed to start containers"})
		return
	}
	c.JSON(200, gin.H{"container": "running"})
}

func RestartContainerHandler(c *gin.Context) {
	container := c.Param("container")
	err := restartContainer(container)
	if err != nil {
		log.Printf("Error restarting container %s: %s", container, err)
		c.JSON(500, gin.H{"error": "Failed to restart containers"})
		return
	}

	c.JSON(200, gin.H{"container": "running"})
}

func StopContainerHandler(c *gin.Context) {
	container := c.Param("container")
	err := stopContainer(container)
	if err != nil {
		log.Printf("Error stopping container %s: %s", container, err)
		c.JSON(500, gin.H{"error": "Failed to stop containers"})
		return
	}
	c.JSON(200, gin.H{"container": "stopped"})
}

func CleanHandler(c *gin.Context) {
	if err := cleanDocker(); err != nil {
		log.Printf("Error cleaning Docker resources: %s", err)
		c.JSON(500, gin.H{"error": "Failed to clean Docker resources"})
		return
	}

	c.JSON(200, gin.H{"message": "Docker resources cleaned successfully"})
}

func cleanDocker() error {
	// Stop and remove all containers
	stopContainersCmd := exec.Command("docker", "stop", "$(docker ps -q)")
	removeContainersCmd := exec.Command("docker", "rm", "$(docker ps -a -q)")
	stopContainersCmd.Run()
	removeContainersCmd.Run()

	// Remove all volumes
	removeVolumesCmd := exec.Command("docker", "volume", "rm", "$(docker volume ls -q)")
	removeVolumesCmd.Run()

	// Remove all images
	removeImagesCmd := exec.Command("docker", "rmi", "$(docker images -q)")
	removeImagesCmd.Run()

	// Prune the Docker system
	pruneCmd := exec.Command("docker", "system", "prune", "-a", "-f")
	pruneCmd.Run()

	return nil
}

func SCPRunHandler(c *gin.Context) {
	script := c.Param("script")
	print(script)
}

func SCPCodeHandler(c *gin.Context) {
	// Parse the multipart form
	err := c.Request.ParseMultipartForm(0) // Set limit to 0 for unlimited file size
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	dirName := c.Request.FormValue("dir-name")
	if dirName == "" {
		c.JSON(400, gin.H{"error": "Missing 'dir-name' parameter"})
		return
	}

	// Get the files from the "file" field in the form
	files, ok := c.Request.MultipartForm.File["file"]
	if !ok || len(files) == 0 {
		c.JSON(400, gin.H{"error": "No files provided"})
		return
	}

	// Specify the destination folder based on the "dir-name" parameter
	destination := dirName
	err = os.MkdirAll(destination, os.ModePerm)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create destination folder"})
		return
	}

	// Iterate through the files and copy them to the destination folder
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to open file"})
			return
		}
		defer src.Close()

		filePath := filepath.Join(destination, file.Filename)
		dst, err := os.Create(filePath)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create the file"})
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to save the file"})
			return
		}
	}

	c.JSON(200, gin.H{"message": "Files uploaded successfully", "dir-name": dirName})
}

func SCPHandler(c *gin.Context) {

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	// Change the destination folder according to your needs
	destination := "./scp/loaded-scripts"
	err = os.MkdirAll(destination, os.ModePerm)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create destination folder"})
		return
	}

	filePath := filepath.Join(destination, header.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create the file"})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save the file"})
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("File %s uploaded successfully", header.Filename)})
}
