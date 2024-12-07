package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"server/adapter/docker"
	"server/adapter/logging"
	"server/domain"
)

type Message struct {
	Message string `json:"message"`
}

type ContainerStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type ContainerMetrics struct {
	Name   string  `json:"name"`
	CPU    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
}

var influxDBClient influxdb2.Client
var org string
var bucket string
var httpRequestCount int64
var mu sync.Mutex

// generateRandomToken generates a random string token of specified length
func generateRandomToken(length int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const tokensFile = "tokens.txt"

	// Create a slice to store the generated token
	token := make([]byte, length)

	// Generate random bytes and convert to token
	for i := range token {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		token[i] = chars[n.Int64()]
	}

	tokenStr := string(token)

	// Open or create the file
	file, err := os.OpenFile(tokensFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the file content
	content, err := os.ReadFile(tokensFile)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	lines := strings.Split(string(content), "\n")

	// Check if the file is empty or not
	if len(lines) == 1 && lines[0] == "" {
		// File exists but is empty, write the token to the first line
		_, err := file.WriteString(tokenStr + "\n")
		if err != nil {
			return "", err
		}
	} else {
		// Append the token to the next empty line or end of the file
		if len(lines) > 0 && lines[len(lines)-1] != "" {
			_, err := file.WriteString("\n" + tokenStr + "\n")
			if err != nil {
				return "", err
			}
		} else {
			_, err := file.WriteString(tokenStr + "\n")
			if err != nil {
				return "", err
			}
		}
	}

	return tokenStr, nil
}

func initInfluxDBClient() {
	influxDBClient = influxdb2.NewClient("http://localhost:8086", "a_secure_admin_token")
	org = "my_org"
	bucket = "my_bucket"
}

func sendMetricsToInfluxDB(metrics []ContainerMetrics) error {
	writeAPI := influxDBClient.WriteAPIBlocking(org, bucket)

	for _, metric := range metrics {
		tags := map[string]string{"container_name": metric.Name}
		fields := map[string]interface{}{
			"cpu":    metric.CPU,
			"memory": metric.Memory,
		}
		p := influxdb2.NewPoint("container_metrics",
			tags,
			fields,
			time.Now())
		err := writeAPI.WritePoint(context.Background(), p)
		if err != nil {
			return fmt.Errorf("failed to write point to InfluxDB: %v", err)
		}
	}

	return nil
}

func sendHttpRequestCountToInfluxDB() error {
	writeAPI := influxDBClient.WriteAPIBlocking(org, bucket)
	tags := map[string]string{"metric": "http_request_count"}
	fields := map[string]interface{}{
		"count": httpRequestCount,
	}
	p := influxdb2.NewPoint("http_requests",
		tags,
		fields,
		time.Now())
	return writeAPI.WritePoint(context.Background(), p)
}

func getContainerMetrics() ([]ContainerMetrics, error) {
	cmd := exec.Command("docker", "stats", "--no-stream", "--format", "\"{{.Name}},{{.CPUPerc}},{{.MemUsage}}\"")
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to execute docker stats: %v", err)
	}

	var metrics []ContainerMetrics
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		line = strings.Trim(line, "\"")
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			continue
		}

		cpuStr := strings.TrimSuffix(parts[1], "%")
		cpu, err := strconv.ParseFloat(cpuStr, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CPU percentage: %v", err)
		}

		memParts := strings.Split(parts[2], " / ")
		if len(memParts) != 2 {
			continue
		}

		memStr := memParts[0]
		memStr = strings.Replace(memStr, "MiB", "", -1)
		memStr = strings.TrimSpace(memStr)
		memory, err := strconv.ParseFloat(memStr, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse memory usage: %v", err)
		}

		metric := ContainerMetrics{
			Name:   parts[0],
			CPU:    cpu,
			Memory: memory,
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	metrics, err := getContainerMetrics()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting container metrics: %v", err), http.StatusInternalServerError)
		return
	}

	err = sendMetricsToInfluxDB(metrics)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending metrics to InfluxDB: %v", err), http.StatusInternalServerError)
		return
	}

	response := Message{Message: "Metrics sent successfully"}
	json.NewEncoder(w).Encode(response)

	err = logging.SendLogToInfluxDB("Metrics collected and sent successfully", "info", influxDBClient, org, bucket)
	if err != nil {
		fmt.Printf("Error logging metrics success: %v\n", err)
	}
}

func checkAuthorizationToken(token string) bool {
	file, err := os.Open("tokens.txt")
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		fmt.Printf("Error opening tokens file: %v\n", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == token {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading tokens file: %v\n", err)
		return false
	}

	return false
}

func genericHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	mu.Lock()
	httpRequestCount++
	mu.Unlock()

	err := sendHttpRequestCountToInfluxDB()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending HTTP request count to InfluxDB: %v", err), http.StatusInternalServerError)
		return
	}

	err = logging.SendLogToInfluxDB(fmt.Sprintf("HTTP request count updated: %d", httpRequestCount), "info", influxDBClient, org, bucket)
	if err != nil {
		fmt.Printf("Error logging HTTP request count: %v\n", err)
	}

	if r.Method == http.MethodPost {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusForbidden)
			return
		}

		// Validate the token
		if !checkAuthorizationToken(authHeader) {
			http.Error(w, "Invalid authorization token", http.StatusForbidden)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var config domain.Config
		if err := json.Unmarshal(body, &config); err != nil {
			http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
			return
		}

		result := docker.StartProcessing(config)

		response := Message{Message: result}
		json.NewEncoder(w).Encode(response)

		err = logging.SendLogToInfluxDB("POST request processed", "info", influxDBClient, org, bucket)
		if err != nil {
			fmt.Printf("Error logging POST request: %v\n", err)
		}
	}

	if r.Method == http.MethodGet {
		statusList, err := docker.GetContainerStatusList()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting container status: %v", err), http.StatusInternalServerError)
			return
		}

		var containers []ContainerStatus
		for name, status := range statusList {
			container := ContainerStatus{
				Name:   name,
				Status: status,
			}
			containers = append(containers, container)
		}

		responseJSON, err := json.Marshal(containers)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error encoding JSON response: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)

		err = logging.SendLogToInfluxDB("GET request processed", "info", influxDBClient, org, bucket)
		if err != nil {
			fmt.Printf("Error logging GET request: %v\n", err)
		}
	}
}

// @title Your API Title
// @version 1.0
// @description Your API Description
// @contact.name Your Name
// @contact.email your.email@example.com
// @host localhost:8081
// @BasePath /
func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "generate-token" {

			token, err := generateRandomToken(60)
			if err != nil {
				fmt.Println("Error generating token:", err)
				return
			}
			fmt.Printf("Your access-token is: %s\n", token)
			return
		} else {
			fmt.Printf("Unknown command argument(s).\n")
			return
		}
	}

	initInfluxDBClient()
	defer influxDBClient.Close()

	port := flag.String("port", "8081", "Port to run the server on")
	flag.Parse()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", genericHandler)
	http.HandleFunc("/metrics", metricsHandler)

	fmt.Printf("Server is listening on port %s...\n", *port)

	err := logging.SendLogToInfluxDB(fmt.Sprintf("Server started on port %s", *port), "info", influxDBClient, org, bucket)
	if err != nil {
		fmt.Printf("Error logging server start: %v\n", err)
	}

	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		err = logging.SendLogToInfluxDB(fmt.Sprintf("Server startup error: %s", err), "error", influxDBClient, org, bucket)
		if err != nil {
			fmt.Printf("Error logging server startup error: %v\n", err)
		}
	}
}
