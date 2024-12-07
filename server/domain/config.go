package domain

// Config struct to unmarshal the JSON configuration
type Config struct {
	AppName        string `json:"app_name"`
	Language       string `json:"language"`
	GitHubRepo     string `json:"github_repo"`
	GitHubUser     string `json:"github_user"`
	GitHubPassword string `json:"github_password"`
	DockerImage    string `json:"docker_image"`
}
