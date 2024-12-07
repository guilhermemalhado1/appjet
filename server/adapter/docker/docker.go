package docker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"server/domain"
	"strings"
)

func StartProcessing(config domain.Config) string {
	switch config.Language {
	case "java":
		generateJavaSetup(config)
		break
	case "python":
		generatePythonSetup(config)
		break
	default:
		return "unsupported yet."
	}

	// Run docker-compose up in the scripts directory
	output, err := runDockerComposeUp()
	if err != nil {
		return fmt.Sprintf("Error running docker-compose up: %v", err)
	}

	return output
}

func runDockerComposeUp() (string, error) {
	// Change working directory to scripts directory
	err := os.Chdir("scripts")
	if err != nil {
		return "", fmt.Errorf("failed to change directory to 'scripts': %v", err)
	}
	defer os.Chdir("..") // Change back to the original directory on function exit

	// Command to run docker-compose up
	cmd := exec.Command("docker-compose", "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start docker-compose up asynchronously
	_ = cmd.Start()

	return "Application online", nil
}

func GetContainerStatusList() (map[string]string, error) {
	// Run docker ps command to get container information in JSON format
	cmd := exec.Command("docker", "ps", "--format", "{{json .}}")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start docker ps command: %v", err)
	}

	// Read output line by line
	scanner := bufio.NewScanner(stdout)
	statusList := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		var container map[string]interface{}
		if err := json.Unmarshal([]byte(line), &container); err != nil {
			return nil, fmt.Errorf("failed to parse line as JSON: %v", err)
		}

		name, ok := container["Names"].(string)
		if !ok {
			return nil, fmt.Errorf("container names assertion failed")
		}

		status := "running"
		if running, ok := container["Running"].(bool); ok && !running {
			status = "not running"
		}

		statusList[strings.TrimLeft(name, "/")] = status
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading docker ps output: %v", err)
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		// Check if the error is due to no containers found, handle it gracefully
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				log.Println("No containers found.")
				return statusList, nil
			}
		}
		return nil, fmt.Errorf("docker ps command failed: %v", err)
	}

	// If statusList is still empty, it means no containers were running
	if len(statusList) == 0 {
		log.Println("No running containers found.")
	}

	log.Println("Returning status list:", statusList)
	return statusList, nil
}

func generateJavaSetup(config domain.Config) string {
	dockerCompose := fmt.Sprintf(`
version: '3.7'
services:
  %s:
    build:
      context: .
      dockerfile: Dockerfile_java
`, config.AppName)

	dockerfile := fmt.Sprintf(`
FROM openjdk:11-jre-slim

# Install git
RUN apt-get update && apt-get install -y git

WORKDIR /app

RUN git clone https://%s:%s@%s /app/repo

CMD ["java", "-jar", "repo/%s.jar"]
`, config.GitHubUser, config.GitHubPassword, config.GitHubRepo, config.AppName)

	// Create scripts directory if it doesn't exist
	err := os.MkdirAll("scripts", os.ModePerm)
	if err != nil {
		return fmt.Sprintf("Failed to create scripts directory: %v", err)
	}

	// Write Docker Compose file
	composeFilePath := filepath.Join("scripts", "docker-compose.yaml")
	if err := writeFile(composeFilePath, dockerCompose); err != nil {
		return fmt.Sprintf("Failed to write Docker Compose file: %v", err)
	}

	// Write Dockerfile
	dockerfileFilePath := filepath.Join("scripts", "Dockerfile_java")
	if err := writeFile(dockerfileFilePath, dockerfile); err != nil {
		return fmt.Sprintf("Failed to write Dockerfile: %v", err)
	}

	return fmt.Sprintf("Docker Compose and Dockerfile generated and placed in 'scripts' folder.")
}

func generatePythonSetup(config domain.Config) string {
	dockerCompose := fmt.Sprintf(`
version: '3.7'
services:
  %s:
    build:
      context: .
      dockerfile: Dockerfile_python
    ports:
      - "8080:8082"
`, config.AppName)

	dockerfile := fmt.Sprintf(`
FROM %s

# Install git and pip
RUN apt-get update && apt-get install -y git python3-pip

# Install Django
RUN pip install django

WORKDIR /app

RUN git clone https://%s:%s@%s .

CMD ["python", "myproject/manage.py", "runserver", "0.0.0.0:8082"]
`, config.DockerImage, config.GitHubUser, config.GitHubPassword, config.GitHubRepo)

	// Create scripts directory if it doesn't exist
	err := os.MkdirAll("scripts", os.ModePerm)
	if err != nil {
		return fmt.Sprintf("Failed to create scripts directory: %v", err)
	}

	// Write Docker Compose file
	composeFilePath := filepath.Join("scripts", "docker-compose.yaml")
	if err := writeFile(composeFilePath, dockerCompose); err != nil {
		return fmt.Sprintf("Failed to write Docker Compose file: %v", err)
	}

	// Write Dockerfile
	dockerfileFilePath := filepath.Join("scripts", "Dockerfile_python")
	if err := writeFile(dockerfileFilePath, dockerfile); err != nil {
		return fmt.Sprintf("Failed to write Dockerfile: %v", err)
	}

	return fmt.Sprintf("Docker Compose and Dockerfile generated and placed in 'scripts' folder.")
}

func writeFile(filePath, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
