package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	// Retrieve command-line arguments
	args := os.Args

	if len(args) < 2 {
		fmt.Println("No command specified.")
		return
	}

	// Define a map to associate each command with its handler function
	commandHandlers := map[string]func([]string){
		"./appjet": handleAppjetCommand,
		"help":     handleHelpCommand,
		/*"logout":   handleLogoutCommand,
		"check-alive": handleCheckAliveCommand,
		"configure":   handleConfigureCommand,
		"inspect":     handleInspectCommand,
		"start":       handleStartCommand,
		"restart":     handleRestartCommand,
		"stop":        handleStopCommand,
		"clean":       handleCleanCommand,
		"scripts":     handleScriptsCommand,
		"code":        handleCodeCommand,
		"scp/run":     handleSCPRunCommand,*/
		"default": handleUnknownCommand,
	}

	// Determine the command and get the corresponding handler
	command := args[1]
	handler, exists := commandHandlers[command]
	if !exists {
		handler = commandHandlers["default"]
	}

	// Call the handler function with the arguments
	handler(args[2:])
}

// Define handler functions for each command

func handleAppjetCommand(arguments []string) {
	fmt.Println("Specify a sub-command. Use './appjet help' for available commands.")
}

func handleHelpCommand(arguments []string) {

	url := "http://localhost:8080/appjet/help"
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

func handleUnknownCommand(arguments []string) {
	fmt.Println("Unknown command. Use './appjet help' for available commands.")
}
