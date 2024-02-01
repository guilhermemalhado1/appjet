package services

import (
	"appjet-server-daemon/app/models"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func GenerateInitSqlIfNeeded(cfg models.Configuration, c *gin.Context) {

	if cfg.Artifact.Application.Language == "python" {
		dbName := cfg.Artifact.Database.Name

		// Use fmt.Sprintf to format the SQL statement with the dynamic database name
		initSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", dbName)

		if err := CreateFile("init.sql", initSql); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Init.sql file"})
			return
		}
	}
}

func GenerateDockerCompose(cfg models.Configuration, c *gin.Context) {

	var volumesInit = ""

	if cfg.Artifact.Application.Language == "python" {
		volumesInit = "volumes:\n      - ./init.sql:/docker-entrypoint-initdb.d/init.sql"
	}

	dockerComposeContent := fmt.Sprintf(`version: '3'
services:
  app:
    container_name: app
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "%d:%d"
    depends_on:
      - database

  database:
    container_name: database
    image: %s
    environment:
      MYSQL_ROOT_PASSWORD: %s
      MYSQL_DATABASE: %s
      MYSQL_ALLOW_EMPTY_PASSWORD: 'yes'
    ports:
      - "%d:%d"
    command: ["mysqld", "--character-set-server=utf8mb4", "--collation-server=utf8mb4_unicode_ci", "--bind-address=0.0.0.0"]
    %s
`,
		cfg.Artifact.Application.Ports.ExternalDocker, cfg.Artifact.Application.Ports.InternalDocker,
		cfg.Artifact.Database.Driver, cfg.Artifact.Database.RootPassword, cfg.Artifact.Database.Name,
		cfg.Artifact.Database.Ports.ExternalDocker, cfg.Artifact.Database.Ports.InternalDocker, volumesInit)

	if err := CreateFile("docker-compose.yml", dockerComposeContent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Docker Compose file"})
		return
	}
}

func GenerateDockerfile(cfg models.Configuration, c *gin.Context) {

	dockerfileContent := ""

	switch cfg.Artifact.Application.Language {
	case "java":
		dockerfileContent = generateJavaDockerFile(cfg)
		if err := CreateFile("Dockerfile", dockerfileContent); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Dockerfile"})
			return
		}
	case "python":
		dockerfileContent = generatePythonDockerFile(cfg)
		if err := CreateFile("Dockerfile", dockerfileContent); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Dockerfile"})
			return
		}
	}
}

func generateJavaDockerFile(cfg models.Configuration) string {
	dockerImage := cfg.Artifact.Application.DockerImage
	target := cfg.Artifact.Application.Artifact.Target
	internalPort := cfg.Artifact.Application.Ports.InternalDocker
	builderName := cfg.Artifact.Application.Builder.Name
	builderDockerImage := cfg.Artifact.Application.Builder.DockerImage

	// Add code checkout logic based on configuration
	var codeCheckoutCmd string
	if cfg.Artifact.CodeCheckout.Git.Enabled {
		codeCheckoutCmd += fmt.Sprintf("ARG GIT_USERNAME=%s\n", cfg.Artifact.CodeCheckout.Git.RepoUser)
		codeCheckoutCmd += fmt.Sprintf("ARG GIT_PASSWORD=%s\n", cfg.Artifact.CodeCheckout.Git.RepoPassword)
		codeCheckoutCmd += "RUN git config --global credential.helper '!f() { echo \"username=${GIT_USERNAME}\"; echo \"password=${GIT_PASSWORD}\"; }; f'\n"
		codeCheckoutCmd += fmt.Sprintf("RUN git clone https://${GIT_USERNAME}:${GIT_PASSWORD}@%s /app || (echo \"Git clone failed\"; exit 1)\n", cfg.Artifact.CodeCheckout.Git.RepoURL)
	}
	if cfg.Artifact.CodeCheckout.SCP.Enabled {
		codeCheckoutCmd += PullCodeFromSCP(cfg.Artifact.CodeCheckout.SCP.Configurations.Folder)
	}

	dockerfile := fmt.Sprintf(`FROM %s as builder

# Install builder dependencies
RUN apt-get update && apt-get install -y git

WORKDIR /app

# Code checkout
%s

COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Build the application
RUN %s install -DskipTests

FROM %s

WORKDIR /app

RUN apt-get update && apt-get install -y netcat

COPY --from=builder /wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

COPY --from=builder %s app.jar

EXPOSE %d

CMD ["/wait-for-it.sh", "database:%d", "--", "java", "-jar", "app.jar"]
`, builderDockerImage, codeCheckoutCmd, getBuilderShortName(builderName), dockerImage, target, internalPort, cfg.Artifact.Database.Ports.InternalDocker)

	return dockerfile
}

func GenerateWaitForItScriptContent() string {
	return `#!/bin/bash
# wait-for-it.sh

set -e

host="$1"
port="$2"
shift 2
cmd="$@"

echo "1 minute before launching the app"
sleep 60

exec $cmd
`
}

func GetDockerContainersState() (map[string]interface{}, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cli.Close()

	containerStates := make(map[string]interface{})

	containerNames := []string{"app", "database"}

	for _, name := range containerNames {
		// Check if the container exists
		exists, err := containerExists(cli, name)
		if err != nil {
			log.Printf("Error checking if container %s exists: %s", name, err)
			return nil, err
		}

		if exists {
			// Check the status using 'docker inspect' command
			statusCmd := exec.Command("docker", "inspect", "--format", "{{.State.Status}}", name)
			statusOutput, err := statusCmd.CombinedOutput()
			if err != nil {
				log.Printf("Error checking container status with 'docker inspect' for %s: %s", name, err)
				return nil, err
			}

			// Trim any leading/trailing whitespaces from the output
			status := strings.TrimSpace(string(statusOutput))
			log.Printf("Container %s status from 'docker inspect': %s", name, status)

			containerStates[name] = map[string]interface{}{
				"status": status == "running",
			}
		} else {
			// Container doesn't exist, report as not running
			containerStates[name] = map[string]interface{}{
				"status": false,
			}
		}
	}

	return containerStates, nil
}

// Check if a container with the given name exists
func containerExists(cli *client.Client, containerName string) (bool, error) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return false, err
	}

	for _, container := range containers {
		for _, containerNameInList := range container.Names {
			if strings.TrimLeft(containerNameInList, "/") == containerName {
				return true, nil
			}
		}
	}

	return false, nil
}

func generatePythonDockerFile(cfg models.Configuration) string {

	var codeCheckoutCmd string
	if cfg.Artifact.CodeCheckout.Git.Enabled {
		codeCheckoutCmd += fmt.Sprintf("ARG GIT_USERNAME=%s\n", cfg.Artifact.CodeCheckout.Git.RepoUser)
		codeCheckoutCmd += fmt.Sprintf("ARG GIT_PASSWORD=%s\n", cfg.Artifact.CodeCheckout.Git.RepoPassword)
		codeCheckoutCmd += "RUN git config --global credential.helper '!f() { echo \"username=${GIT_USERNAME}\"; echo \"password=${GIT_PASSWORD}\"; }; f'\n"
		codeCheckoutCmd += fmt.Sprintf("RUN git clone https://${GIT_USERNAME}:${GIT_PASSWORD}@%s /app_builder || (echo \"Git clone failed\"; exit 1)\n", cfg.Artifact.CodeCheckout.Git.RepoURL)
	}
	if cfg.Artifact.CodeCheckout.SCP.Enabled {
		codeCheckoutCmd += PullCodeFromSCP(cfg.Artifact.CodeCheckout.SCP.Configurations.Folder)
	}

	return fmt.Sprintf(`# Use an official Python runtime as a parent image
FROM %s AS builder

# Set environment variables
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

# Install system dependencies
RUN apt-get update \
    && apt-get install -y --no-install-recommends git pkg-config libmariadb-dev-compat build-essential \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /app_builder

# Code checkout
%s

COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Print the contents of the cloned directory (for debugging)
RUN ls -la /app_builder

# Set the working directory again (just to be explicit)
WORKDIR /app_builder

# Create and activate a virtual environment
RUN python -m venv venv
ENV PATH="/app_builder/venv/bin:$PATH"

# Install dependencies from requirements.txt within the virtual environment
RUN . /app_builder/venv/bin/activate && pip install --upgrade pip \
    && pip install -r /app_builder/requirements.txt

# Second stage for the final image
FROM %s

# Set environment variables
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

# Set the working directory
WORKDIR /app

# Copy files from the builder stage
COPY --from=builder /app_builder /app

COPY --from=builder /wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Expose the port the app runs on
EXPOSE %d

# Set the working directory for the CMD
WORKDIR /app

# Create the database, run migrations, and start the server
CMD /wait-for-it.sh database:%d -- /app/venv/bin/python manage.py makemigrations && /app/venv/bin/python manage.py migrate && /app/venv/bin/python manage.py runserver 0.0.0.0:%d
`, cfg.Artifact.Application.DockerImage, codeCheckoutCmd,
		cfg.Artifact.Application.DockerImage, cfg.Artifact.Application.Ports.InternalDocker, cfg.Artifact.Database.Ports.InternalDocker, cfg.Artifact.Application.Ports.InternalDocker)
}

func getBuilderShortName(builderName string) string {
	switch builderName {
	case "maven":
		return "mvn"
	case "gradle":
		return "gradle"
	default:
		return ""
	}
}
