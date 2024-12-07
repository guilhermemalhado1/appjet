package main

import (
	"client/handler"
	"fmt"
	"os"

	"client/ui"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected one of the following subcommands: [deploy, ui]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "deploy":
		if len(os.Args) > 2 {
			handler.RunCommand(os.Args[2])
			os.Exit(1)
		}
		fmt.Printf("Missing token parameter in command.\n")
		os.Exit(1)
	case "status":
		if len(os.Args) > 2 {
			handler.GetStatus(os.Args[2])
			os.Exit(1)
		}
		handler.GetStatus("")
		os.Exit(1)
	case "ui":
		ui.RunApp()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
