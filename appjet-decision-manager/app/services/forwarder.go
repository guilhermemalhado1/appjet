package services

import (
	"appjet-decision-manager/app/models"
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func ForwardConfigToDaemon(config *models.Configuration, url string) (resp *http.Response, err error) {
	// Convert config to JSON
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	// Make HTTP POST request to the daemon
	response, err := http.Post(url, "application/json", bytes.NewBuffer(configJSON))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return response, err
}

func ForwardStartToDaemon(url string) (resp *http.Response, err error) {

	// Make HTTP POST request to the daemon
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return response, err
}

func ForwardRestartToDaemon(url string) (resp *http.Response, err error) {

	// Make HTTP POST request to the daemon
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return response, err
}

func ForwardStopToDaemon(url string) (resp *http.Response, err error) {

	// Make HTTP POST request to the daemon
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return response, err
}

func ForwardCheckAliveToDaemon(url string) (resp *http.Response, err error) {

	// Make HTTP POST request to the daemon
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return response, err
}

func ForwardInspectToDaemon(url string) (resp *http.Response, err error) {

	// Make HTTP POST request to the daemon
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return response, err
}

func ForwardCleanToDaemon(url string) (resp *http.Response, err error) {

	// Make HTTP POST request to the daemon
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return response, err
}

func ForwardSCPRunToDaemon(url string) (resp *http.Response, err error) {

	// Make HTTP POST request to the daemon
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return response, err
}

func ForwardSCPToDaemon(scriptBytes multipart.File, filename string, url string) (resp *http.Response, err error) {
	// Create a new buffer to store the multipart request body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a form file field for the scriptBytes
	part, err := writer.CreateFormFile("file", filename) // Change "filename" to an appropriate name
	if err != nil {
		return nil, err
	}

	// Copy the content of scriptBytes to the form file field
	_, err = io.Copy(part, scriptBytes)
	if err != nil {
		return nil, err
	}

	// Close the multipart writer before making the request
	writer.Close()

	// Make HTTP POST request to the daemon
	response, err := http.Post(url, writer.FormDataContentType(), body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return response, err
}

func ForwardSCPCodeToDaemon(codeDirectory string, url string) (resp *http.Response, err error) {
	// Create a new buffer to store the multipart request body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Walk through the specified directory and its subdirectories
	err = filepath.Walk(codeDirectory, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Create a form file field for each file
		part, err := writer.CreateFormFile("files", filePath)
		if err != nil {
			return err
		}

		// Open the file
		fileContent, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer fileContent.Close()

		// Copy the content of the file to the form file field
		_, err = io.Copy(part, fileContent)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Close the multipart writer before making the request
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Make HTTP POST request to the SCPCodeHandler endpoint
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return response, nil
}
