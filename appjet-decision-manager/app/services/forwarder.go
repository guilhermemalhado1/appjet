package services

import (
	"appjet-decision-manager/app/models"
	"bytes"
	"encoding/json"
	"net/http"
)

func ForwardConfigToDaemon(config models.Configuration, url string) (resp *http.Response, err error) {
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
