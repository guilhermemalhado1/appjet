package services

import (
	"appjet-cli/app/models"
	"encoding/json"
	"io/ioutil"
)

func GetConfiguration() (models.Configuration, error) {
	var config models.Configuration

	// Read the contents of config.json
	configJSON, err := ioutil.ReadFile("config.json")
	if err != nil {
		return config, err
	}

	// Unmarshal the JSON into the config variable
	err = json.Unmarshal(configJSON, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
