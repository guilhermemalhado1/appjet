package ui

import (
	"client/handler"
	"encoding/json"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	myApp            fyne.App
	myWindow         fyne.Window
	fileContentText  *widget.Label
	fileNameLabel    *widget.Label
	startButton      *widget.Button
	uploadButton     *widget.Button
	tokenEntry       *widget.Entry
	selectedFilePath string
	targetURLs       []string
	statusContent    *widget.Label
	statusTab        *container.TabItem
)

// RunApp initializes and runs the GUI application
func RunApp() {
	myApp = app.New()
	myWindow = myApp.NewWindow("Appjet Client UI")

	// Set the size of the window to 800x600 pixels
	myWindow.Resize(fyne.NewSize(800, 600))

	// Initialize widgets
	fileContentText = widget.NewLabel("")
	fileNameLabel = widget.NewLabel("No file selected")
	startButton = widget.NewButtonWithIcon("Start Processing", theme.ContentAddIcon(), func() {
		token := tokenEntry.Text
		if token == "" {
			dialog.ShowError(errors.New("Access Token cannot be empty"), myWindow)
			return
		}
		saveAndProcessFile(token)
	})
	startButton.Disable() // Start button initially disabled

	uploadButton = widget.NewButtonWithIcon("Upload JSON File", theme.FileIcon(), func() {
		openFile()
	})

	tokenEntry = widget.NewEntry()
	tokenEntry.SetPlaceHolder("Enter Access Token")

	// Enable the startButton only if the tokenEntry is not empty
	tokenEntry.OnChanged = func(text string) {
		if text != "" {
			startButton.Enable()
		} else {
			startButton.Disable()
		}
	}

	// Create a "Status" tab content
	statusContent = widget.NewLabel("Status will be displayed here")
	statusTab = container.NewTabItem("Status", container.NewScroll(statusContent))

	// Create tabs for "Upload" and "Status"
	tabs := container.NewAppTabs(
		container.NewTabItem("Upload", createUploadTab()),
		statusTab,
	)

	// Set the content of the main window
	myWindow.SetContent(container.NewBorder(nil, nil, nil, nil, tabs))

	// Handle tab changes
	tabs.OnChanged = func(tab *container.TabItem) {
		if tab == statusTab {
			updateStatusTab()
		}
	}

	// Show and run the application
	myWindow.ShowAndRun()
}

// createUploadTab creates the content for the Upload tab
func createUploadTab() fyne.CanvasObject {
	// Create styled containers
	buttonsContainer := container.New(layout.NewVBoxLayout(),
		container.New(layout.NewHBoxLayout(),
			layout.NewSpacer(),
			uploadButton,
			layout.NewSpacer(),
		),
		container.New(layout.NewHBoxLayout(),
			layout.NewSpacer(),
			startButton,
			layout.NewSpacer(),
		),
	)

	// Create a vertical box container for the content with padding
	content := container.NewVBox(
		widget.NewLabelWithStyle("Upload and process the config.json file.", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(" "),
		tokenEntry,
		widget.NewLabel(" "),
		buttonsContainer,
		widget.NewLabel(" "),
		fileNameLabel,
		widget.NewLabelWithStyle("File Content:", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		fileContentText,
	)

	// Center the content vertically and horizontally with padding
	return container.New(layout.NewVBoxLayout(),
		layout.NewSpacer(),
		content,
		layout.NewSpacer(),
	)
}

// openFile opens a file dialog to upload a JSON file
func openFile() {
	fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		if reader == nil {
			// Dialog was cancelled
			return
		}
		// Check if the file has a .json extension
		if !strings.HasSuffix(strings.ToLower(reader.URI().String()), ".json") {
			dialog.ShowError(errors.New("only JSON files are allowed"), myWindow)
			return
		}

		// Update file name label
		fileNameLabel.SetText(fmt.Sprintf("Selected file: %s", reader.URI().Name()))

		// Read file content
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

		// Format JSON content
		var jsonData map[string]interface{}
		if err := json.Unmarshal(data, &jsonData); err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		prettyData, err := json.MarshalIndent(jsonData, "", "  ")
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

		// Extract target URLs from JSON
		if urls, ok := jsonData["target_urls"].([]interface{}); ok {
			targetURLs = make([]string, len(urls))
			for i, url := range urls {
				if urlMap, ok := url.(map[string]interface{}); ok {
					if urlStr, ok := urlMap["url"].(string); ok {
						targetURLs[i] = urlStr
					}
				}
			}
		}

		// Set formatted content
		fileContentText.SetText(string(prettyData))

		// Enable the startButton
		startButton.Enable()

		// Save the file path
		selectedFilePath = reader.URI().Path()

		// Close the reader
		if err := reader.Close(); err != nil {
			dialog.ShowError(err, myWindow)
		}
	}, myWindow)

	// Show the file dialog
	fileDialog.Show()
}

// saveAndProcessFile saves the uploaded file and calls the handler method
func saveAndProcessFile(token string) {
	// Define the target file path in the current directory
	targetFilePath := filepath.Join(".", "config.json")

	// Read the selected file content
	data, err := ioutil.ReadFile(selectedFilePath)
	if err != nil {
		dialog.ShowError(err, myWindow)
		return
	}

	// Save the content to the target file path
	err = ioutil.WriteFile(targetFilePath, data, 0644)
	if err != nil {
		dialog.ShowError(err, myWindow)
		return
	}

	// Create a custom loading dialog with improved styling
	loadingDialog := dialog.NewCustom("Processing", "Close", container.NewVBox(
		widget.NewLabel("Processing..."),
		widget.NewProgressBarInfinite(),
	), myWindow)

	// Show the loading dialog
	loadingDialog.Show()

	// Call the handler method with the specified arguments
	go func() {
		success := handler.RunCommand(token)

		// Hide the loading dialog and update status on the main thread
		fyne.CurrentApp().Driver().CanvasForObject(myWindow.Content()).Refresh(myWindow.Content())

		// Show success or error message
		if success {
			updateStatusTab()
			dialog.ShowInformation("Success", "Processing completed successfully. Check the Status tab.", myWindow)
		} else {
			dialog.ShowError(errors.New("Processing failed. Check if the Access Token is still valid"), myWindow)
		}
	}()
}

// updateStatusTab fetches and displays JSON array responses in the "Status" tab
func updateStatusTab() {
	// Set the initial message to indicate fetching is in progress
	statusContent.SetText("Fetching status responses...\n")

	go func() {
		// Call handler.GetStatusJSON and get the status response
		statusResponse := handler.GetStatusJSON() // Assuming this method returns a formatted JSON string

		// Update the statusContent label on the main thread
		fyne.CurrentApp().Driver().CanvasForObject(myWindow.Content()).Refresh(myWindow.Content())
		statusContent.SetText(statusResponse)
	}()
}

// openURLInBrowser opens the provided URL in the default web browser
func openURLInBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		dialog.ShowError(errors.New("unsupported operating system"), myWindow)
		return
	}
	if err := cmd.Start(); err != nil {
		dialog.ShowError(err, myWindow)
	}
}
