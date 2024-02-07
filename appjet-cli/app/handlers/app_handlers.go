package handlers

import (
	"appjet-cli/app/models"
	"appjet-cli/app/services"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func HandleHelpCommand(arguments []string, config models.Configuration) {

	url := config.IdentityProvider.ServerURL + "/appjet/help"
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making HTTP GET request:", err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Format and print the JSON
	var formattedJSON bytes.Buffer
	err = json.Indent(&formattedJSON, body, "", "    ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}

	fmt.Println("HTTP GET request to", url, "completed. Response:")
	fmt.Println(formattedJSON.String())
}

func HandleLoginCommand(arguments []string, config models.Configuration) {
	var username, password string

	// Ask for the username
	fmt.Print("Enter username: ")
	fmt.Scanln(&username)

	// Ask for the password with masking
	fmt.Print("Enter password: ")
	fmt.Scanln(&password)

	// Prepare the request body
	requestBody, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		fmt.Println("Error marshalling request body:", err)
		return
	}

	// Make HTTP POST request
	url := config.IdentityProvider.ServerURL + "/appjet/login"
	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error making HTTP POST request:", err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Check the response status code
	if response.StatusCode == http.StatusOK {
		err := services.EncryptAndSaveToken(string(body))
		if err != nil {
			fmt.Println("Error generating token security file:", err)
			return
		}
		fmt.Println("Login successful. ")
		// Proceed with the logic for successful login
	} else {
		fmt.Println("Login failed. Status code:", response.StatusCode)
		// Handle failed login attempt
	}
}

func HandleLogoutCommand(arguments []string, config models.Configuration) {

	token, _ := services.DecryptToken()

	// Construct the logout URL with the token
	url := fmt.Sprintf("%s/appjet/logout/%s", config.IdentityProvider.ServerURL, token)

	// Make HTTP GET request
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making HTTP GET request:", err)
		return
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode == http.StatusOK {
		fmt.Println("Logout successful")
		// Proceed with the logic for successful logout
	} else {
		fmt.Println("Logout failed. Status code:", response.StatusCode)
		// Handle failed logout attempt
	}
}

func HandleCheckAliveCommand(arguments []string, config models.Configuration) {
	token, _ := services.DecryptToken()

	var url string
	switch len(arguments) {
	case 0:
		url = config.IdentityProvider.ServerURL + "/appjet/check-alive"
	case 1:
		url = fmt.Sprintf("%s/appjet/check-alive/%s", config.IdentityProvider.ServerURL, arguments[0])
	case 2:
		url = fmt.Sprintf("%s/appjet/check-alive/%s/%s", config.IdentityProvider.ServerURL, arguments[0], arguments[1])
	default:
		fmt.Println("Invalid number of arguments")
		return
	}

	fmt.Println("Token:", token)
	fmt.Println("URL:", url)

	makeGETRequest(url)

}

func HandleConfigureCommand(arguments []string, config models.Configuration) {
	token, err := services.DecryptToken()
	if err != nil {
		fmt.Println("Error decrypting token:", err)
		return
	}

	// Convert the config struct to JSON
	configJSON, err := json.Marshal(config)
	if err != nil {
		fmt.Println("Error marshaling configuration to JSON:", err)
		return
	}

	var url string
	switch len(arguments) {
	case 0:
		url = config.IdentityProvider.ServerURL + "/appjet/configure"
	case 1:
		url = fmt.Sprintf("%s/appjet/configure/%s", config.IdentityProvider.ServerURL, arguments[0])
	case 2:
		url = fmt.Sprintf("%s/appjet/configure/%s/%s", config.IdentityProvider.ServerURL, arguments[0], arguments[1])
	default:
		fmt.Println("Invalid number of arguments")
		return
	}

	// Prepare the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(configJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	// Set the Authorization header with the decrypted token
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	// Print the response status
	fmt.Println("Response Status:", resp.Status)
}

func HandleInspectCommand(arguments []string, config models.Configuration) {
	// Your existing code here
	token, _ := services.DecryptToken()

	var url string
	switch len(arguments) {
	case 0:
		url = config.IdentityProvider.ServerURL + "/appjet/inspect"
	case 1:
		url = fmt.Sprintf("%s/appjet/inspect/%s", config.IdentityProvider.ServerURL, arguments[0])
	case 2:
		url = fmt.Sprintf("%s/appjet/inspect/%s/%s", config.IdentityProvider.ServerURL, arguments[0], arguments[1])
	case 3:
		url = fmt.Sprintf("%s/appjet/inspect/%s/%s/%s", config.IdentityProvider.ServerURL, arguments[0], arguments[1], arguments[2])
	default:
		fmt.Println("Invalid number of arguments")
		return
	}

	fmt.Println("Token:", token)
	fmt.Println("URL:", url)

	makeGETRequest(url)
}

func HandleStartCommand(arguments []string, config models.Configuration) {
	// Your existing code here
	token, _ := services.DecryptToken()

	var url string
	switch len(arguments) {
	case 0:
		url = config.IdentityProvider.ServerURL + "/appjet/start"
	case 1:
		url = fmt.Sprintf("%s/appjet/start/%s", config.IdentityProvider.ServerURL, arguments[0])
	case 2:
		url = fmt.Sprintf("%s/appjet/start/%s/%s", config.IdentityProvider.ServerURL, arguments[0], arguments[1])
	case 3:
		url = fmt.Sprintf("%s/appjet/start/%s/%s/%s", config.IdentityProvider.ServerURL, arguments[0], arguments[1], arguments[2])
	default:
		fmt.Println("Invalid number of arguments")
		return
	}

	fmt.Println("Token:", token)
	fmt.Println("URL:", url)

	makeGETRequest(url)
}

func HandleRestartCommand(arguments []string, config models.Configuration) {
	// Your existing code here
	token, _ := services.DecryptToken()

	var url string
	switch len(arguments) {
	case 0:
		url = config.IdentityProvider.ServerURL + "/appjet/restart"
	case 1:
		url = fmt.Sprintf("%s/appjet/restart/%s", config.IdentityProvider.ServerURL, arguments[0])
	case 2:
		url = fmt.Sprintf("%s/appjet/restart/%s/%s", config.IdentityProvider.ServerURL, arguments[0], arguments[1])
	case 3:
		url = fmt.Sprintf("%s/appjet/restart/%s/%s/%s", config.IdentityProvider.ServerURL, arguments[0], arguments[1], arguments[2])
	default:
		fmt.Println("Invalid number of arguments")
		return
	}

	fmt.Println("Token:", token)
	fmt.Println("URL:", url)

	makeGETRequest(url)
}

func HandleStopCommand(arguments []string, config models.Configuration) {
	// Your existing code here
	token, _ := services.DecryptToken()

	var url string
	switch len(arguments) {
	case 0:
		url = config.IdentityProvider.ServerURL + "/appjet/stop"
	case 1:
		url = fmt.Sprintf("%s/appjet/stop/%s", config.IdentityProvider.ServerURL, arguments[0])
	case 2:
		url = fmt.Sprintf("%s/appjet/stop/%s/%s", config.IdentityProvider.ServerURL, arguments[0], arguments[1])
	case 3:
		url = fmt.Sprintf("%s/appjet/stop/%s/%s/%s", config.IdentityProvider.ServerURL, arguments[0], arguments[1], arguments[2])
	default:
		fmt.Println("Invalid number of arguments")
		return
	}

	fmt.Println("Token:", token)
	fmt.Println("URL:", url)

	makeGETRequest(url)
}

func HandleCleanCommand(arguments []string, config models.Configuration) {
	// Your existing code here
	token, _ := services.DecryptToken()

	var url string
	switch len(arguments) {
	case 0:
		url = config.IdentityProvider.ServerURL + "/appjet/clean"
	case 1:
		url = fmt.Sprintf("%s/appjet/clean/%s", config.IdentityProvider.ServerURL, arguments[0])
	case 2:
		url = fmt.Sprintf("%s/appjet/clean/%s/%s", config.IdentityProvider.ServerURL, arguments[0], arguments[1])
	case 3:
		url = fmt.Sprintf("%s/appjet/clean/%s/%s/%s", config.IdentityProvider.ServerURL, arguments[0], arguments[1], arguments[2])
	default:
		fmt.Println("Invalid number of arguments")
		return
	}

	fmt.Println("Token:", token)
	fmt.Println("URL:", url)

	makeGETRequest(url)
}

func HandleScriptsCommand(arguments []string, config models.Configuration) {
	token, err := services.DecryptToken()
	if err != nil {
		fmt.Println("Error decrypting token:", err)
		return
	}

	// Open the directory containing scripts
	dir := "./scripts" // Assuming "scripts" directory is in the current folder
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// Prompt the user to select the script
	fmt.Println("Select the script to upload:")
	for i, file := range files {
		fmt.Printf("%d. %s\n", i+1, file.Name())
	}

	// Read user input
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the number of the script: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Convert input to integer
	scriptNumber, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || scriptNumber < 1 || scriptNumber > len(files) {
		fmt.Println("Invalid script number.")
		return
	}

	// Get the selected script file
	selectedFile := files[scriptNumber-1]

	// Open the selected script file
	filePath := filepath.Join(dir, selectedFile.Name())
	scriptFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer scriptFile.Close()

	// Read the file content
	fileBytes, err := ioutil.ReadAll(scriptFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Create a new buffer to store the multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add the file content to the multipart form data
	fileWriter, err := writer.CreateFormFile("file", selectedFile.Name())
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}
	_, err = fileWriter.Write(fileBytes)
	if err != nil {
		fmt.Println("Error writing file to form file:", err)
		return
	}

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing multipart writer:", err)
		return
	}

	// Construct the URL based on the provided arguments
	var url string
	switch len(arguments) {
	case 0:
		url = config.IdentityProvider.ServerURL + "/appjet/scripts"
	case 1:
		url = fmt.Sprintf("%s/appjet/scripts/%s", config.IdentityProvider.ServerURL, arguments[0])
	case 2:
		url = fmt.Sprintf("%s/appjet/scripts/%s/%s", config.IdentityProvider.ServerURL, arguments[0], arguments[1])
	default:
		fmt.Println("Invalid number of arguments")
		return
	}

	// Prepare the request
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	// Set the content type for multipart form data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Set the Authorization header with the decrypted token
	req.Header.Set("Authorization", token)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	// Print the response status
	fmt.Println("Response Status:", resp.Status)
}

func HandleCodeCommand(arguments []string, config models.Configuration) {
	token, _ := services.DecryptToken()

	var url string
	switch len(arguments) {
	case 0:
		url = config.IdentityProvider.ServerURL + "/appjet/code"
	case 1:
		url = fmt.Sprintf("%s/appjet/code/%s", config.IdentityProvider.ServerURL, arguments[0])
	case 2:
		url = fmt.Sprintf("%s/appjet/code/%s/%s", config.IdentityProvider.ServerURL, arguments[0], arguments[1])
	default:
		fmt.Println("Invalid number of arguments")
		return
	}

	fmt.Println("Token:", token)
	fmt.Println("URL:", url)

	// Send files via SCP
	err := sendFilesViaSCP(config.IdentityProvider.ServerURL, config.IdentityProvider.ServerUsername, config.IdentityProvider.ServerPassword, arguments)
	if err != nil {
		fmt.Println("Error sending files via SCP:", err)
		return
	}

	// Make the POST request
	resp, err := makePOSTRequest(url)
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}
	defer resp.Body.Close()

	// Print the response status
	fmt.Println("Response Status:", resp.Status)
}

func HandleSCPRunCommand(arguments []string, config models.Configuration) {
	// Your existing code here
	token, _ := services.DecryptToken()

	var url string
	switch len(arguments) {
	case 0:
		url = config.IdentityProvider.ServerURL + "/appjet/scp/run"
	case 1:
		url = fmt.Sprintf("%s/appjet/scp/run/%s", config.IdentityProvider.ServerURL, arguments[0])
	case 2:
		url = fmt.Sprintf("%s/appjet/scp/run/%s/%s", config.IdentityProvider.ServerURL, arguments[0], arguments[1])
	case 3:
		url = fmt.Sprintf("%s/appjet/scp/run/%s/%s/%s", config.IdentityProvider.ServerURL, arguments[0], arguments[1], arguments[2])
	default:
		fmt.Println("Invalid number of arguments")
		return
	}

	fmt.Println("Token:", token)
	fmt.Println("URL:", url)

	makeGETRequest(url)
}

func HandleUnknownCommand(arguments []string, config models.Configuration) {
	print("Unknown command.")
}

func makeGETRequest(url string) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making HTTP GET request:", err)
		return
	}
	defer response.Body.Close()

	// Read the response body, if needed
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Process the response as needed
	fmt.Println("Response:", string(body))
}

func makePOSTRequest(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

func sendFilesViaSCP(serverURL, username, password string, arguments []string) error {
	// Connect to the server via SSH
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshClient, err := ssh.Dial("tcp", serverURL+":22", sshConfig)
	if err != nil {
		return err
	}
	defer sshClient.Close()

	// Create an SFTP session
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	// Prepare the source and target directories
	sourceDir := "./code"
	targetDir := "/code"
	if len(arguments) > 0 {
		targetDir = filepath.Join(targetDir, arguments[0]) // Append the first argument as a subdirectory on the server
	}

	// Walk through the source directory and transfer files to the target directory
	err = filepath.Walk(sourceDir, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if fileInfo.IsDir() {
			return nil
		}

		// Open the local file
		localFile, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}

		// Create the remote file
		remoteFile, err := sftpClient.Create(filepath.Join(targetDir, fileInfo.Name()))
		if err != nil {
			return err
		}
		defer remoteFile.Close()

		// Write the contents of the local file to the remote file
		_, err = remoteFile.Write(localFile)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
