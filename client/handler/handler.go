package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Define structs to match the configuration and post data
type TargetURL struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type Config struct {
	AppName        string      `json:"app_name"`
	Language       string      `json:"language"`
	GithubRepo     string      `json:"github_repo"`
	GithubUser     string      `json:"github_user"`
	GithubPassword string      `json:"github_password"`
	DockerImage    string      `json:"docker_image"`
	TargetURLs     []TargetURL `json:"target_urls"`
}

type PostData struct {
	AppName        string `json:"app_name"`
	Language       string `json:"language"`
	GithubRepo     string `json:"github_repo"`
	GithubUser     string `json:"github_user"`
	GithubPassword string `json:"github_password"`
	DockerImage    string `json:"docker_image"`
}

func GetStatusJSON() string {
	var results []map[string]interface{}

	// Read the content of the config.json file
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return fmt.Sprintf("Error reading config.json: %s\n", err)
	}

	// Unmarshal the JSON data into a Config struct
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return fmt.Sprintf("Error unmarshalling config.json: %s\n", err)
	}

	// Fetch status from each URL in the target_urls array
	for _, target := range config.TargetURLs {
		url := target.URL

		resp, err := http.Get(url)
		if err != nil {
			results = append(results, map[string]interface{}{
				"url":   url,
				"error": fmt.Sprintf("Error fetching status from %s: %s", url, err),
			})
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			results = append(results, map[string]interface{}{
				"url":   url,
				"error": fmt.Sprintf("Error reading response body from %s: %s", url, err),
			})
			continue
		}

		return "Status of Server: " + target.URL + "\n\n" + string(body)

	}

	// Marshal the results into a JSON array
	resultJSON, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error marshalling results to JSON: %s\n", err)
	}

	return string(resultJSON)
}

// GetStatus sends HTTP requests based on the serverName
func GetStatus(serverName string) bool {
	// Read the content of the config.json file from the current directory
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error reading config.json:", err)
		return false
	}

	// Unmarshal the JSON data into a Config struct
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error unmarshalling config.json:", err)
		return false
	}

	// Create an HTTP client
	client := &http.Client{}

	if serverName == "" {
		// Send HTTP request to all servers.URL
		for _, target := range config.TargetURLs {
			req, err := http.NewRequest("GET", target.URL, nil)
			if err != nil {
				fmt.Printf("Error creating request for %s: %s\n", target.URL, err)
				return false
			}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Error making GET request to %s: %s\n", target.URL, err)
				return false
			}
			defer resp.Body.Close()

			// Check the status code
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Received non-200 response code from %s: %d\n", target.URL, resp.StatusCode)
				return false
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error reading response body from %s: %s\n", target.URL, err)
				return false
			}
			fmt.Printf("Status from %s:\n%s\n", target.URL, body)
		}
	} else {
		// Send HTTP request to the server with the matching name
		var found bool
		for _, target := range config.TargetURLs {
			if target.Name == serverName {
				req, err := http.NewRequest("GET", target.URL, nil)
				if err != nil {
					fmt.Printf("Error creating request for %s: %s\n", target.URL, err)
					return false
				}
				resp, err := client.Do(req)
				if err != nil {
					fmt.Printf("Error making GET request to %s: %s\n", target.URL, err)
					return false
				}
				defer resp.Body.Close()

				// Check the status code
				if resp.StatusCode != http.StatusOK {
					fmt.Printf("Received non-200 response code from %s: %d\n", target.URL, resp.StatusCode)
					return false
				}

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("Error reading response body from %s: %s\n", target.URL, err)
					return false
				}

				fmt.Printf("Status from %s:\n%s\n", target.URL, body)
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("No server found with the name %s\n", serverName)
			return false
		}
	}

	return true
}

// RunCommand handles the main execution logic
func RunCommand(token string) bool {
	// Read the content of the config.json file from the current directory
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error reading config.json:", err)
		return false
	}

	// Unmarshal the JSON data into a Config struct
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error unmarshalling config.json:", err)
		return false
	}

	// Create PostData struct without TargetURLs
	postData := PostData{
		AppName:        config.AppName,
		Language:       config.Language,
		GithubRepo:     config.GithubRepo,
		GithubUser:     config.GithubUser,
		GithubPassword: config.GithubPassword,
		DockerImage:    config.DockerImage,
	}

	// Marshal the PostData struct into JSON
	postDataJSON, err := json.Marshal(postData)
	if err != nil {
		fmt.Println("Error marshalling PostData:", err)
		return false
	}

	// Send a POST request to each URL in the target_urls array
	for _, elem := range config.TargetURLs {
		req, err := http.NewRequest("POST", elem.URL, bytes.NewReader(postDataJSON))
		if err != nil {
			fmt.Printf("Error creating request for %s: %s\n", elem.URL, err)
			return false
		}
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making POST request to %s: %s\n", elem.URL, err)
			return false
		}
		defer resp.Body.Close()

		// Check the status code
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Received non-200 response code from %s: %d\n", elem.URL, resp.StatusCode)
			return false
		}
	}

	// Wait for 10 seconds
	time.Sleep(10 * time.Second)

	// Print URLs for console output
	fmt.Println("Processing completed successfully. The website is running at the following URLs:")
	for _, elem := range config.TargetURLs {
		elem.URL = elem.URL[:len(elem.URL)-4] + "8080"
		fmt.Printf("%s\n", elem.URL)
	}

	return true
}
