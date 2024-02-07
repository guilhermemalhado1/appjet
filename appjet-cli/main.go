package main

import (
	handlers "appjet-cli/app/handlers"
	"appjet-cli/app/models"
	services "appjet-cli/app/services"
	"fmt"
	"os"
)

func main() {
	// Retrieve command-line arguments
	args := os.Args

	if len(args) < 2 {
		fmt.Println("No command specified.")
		return
	}

	commandHandlers := map[string]func([]string, models.Configuration){
		"help":        handlers.HandleHelpCommand,
		"login":       handlers.HandleLoginCommand,
		"logout":      handlers.HandleLogoutCommand,
		"check-alive": handlers.HandleCheckAliveCommand,
		"configure":   handlers.HandleConfigureCommand,
		"inspect":     handlers.HandleInspectCommand,
		"start":       handlers.HandleStartCommand,
		"restart":     handlers.HandleRestartCommand,
		"stop":        handlers.HandleStopCommand,
		"clean":       handlers.HandleCleanCommand,
		"scripts":     handlers.HandleScriptsCommand,
		"code":        handlers.HandleCodeCommand,
		"scp/run":     handlers.HandleSCPRunCommand,
		"default":     handlers.HandleUnknownCommand,
	}

	config, _ := services.GetConfiguration()

	command := args[1]
	handler, exists := commandHandlers[command]
	if !exists {
		handler = commandHandlers["default"]
	}

	handler(args[2:], config)
}
