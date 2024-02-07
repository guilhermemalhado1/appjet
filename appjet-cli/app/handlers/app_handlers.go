package handlers

import (
	"appjet-cli/app/models"
	"appjet-cli/app/services"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	// Your existing code here
	token, _ := services.DecryptToken()
}

func HandleInspectCommand(arguments []string, config models.Configuration) {
	// Your existing code here
	token, _ := services.DecryptToken()
}

func HandleStartCommand(arguments []string, config models.Configuration) {
	// Your existing code here
	token, _ := services.DecryptToken()
}

func HandleRestartCommand(arguments []string, config models.Configuration) {
	// Your existing code here
	token, _ := services.DecryptToken()
}

func HandleStopCommand(arguments []string, config models.Configuration) {
	// Your existing code here
	token, _ := services.DecryptToken()
}

func HandleCleanCommand(arguments []string, config models.Configuration) {
	// Your existing code here
	token, _ := services.DecryptToken()
}

func HandleScriptsCommand(arguments []string, config models.Configuration) {
	// Your existing code here
	token, _ := services.DecryptToken()
}

func HandleCodeCommand(arguments []string, config models.Configuration) {
	// Your existing code here
	token, _ := services.DecryptToken()
}

func HandleSCPRunCommand(arguments []string, config models.Configuration) {
	// Your existing code here
	token, _ := services.DecryptToken()
}

func HandleUnknownCommand(arguments []string, config models.Configuration) {
	// Your existing code here
	token, _ := services.DecryptToken()
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
