package models

type Configuration struct {
	IdentityProvider struct {
		ServerURL      string `json:"server-url"`
		ServerUsername string `json:"server-username"`
		ServerPassword string `json:"server-password"`
	} `json:"identity-provider"`

	Clusters []struct {
		Name    string `json:"name"`
		Servers []struct {
			Name          string `json:"name"`
			IP            string `json:"ip"`
			Port          int    `json:"port"`
			User          string `json:"user"`
			Password      string `json:"password"`
			DeployDetails struct {
				Folder string `json:"folder"`
			} `json:"deploy-details"`
		} `json:"servers"`
	} `json:"clusters"`

	Artifact struct {
		Application struct {
			Language  string `json:"language"`
			Framework string `json:"framework"`
			Artifact  struct {
				Target string `json:"target"`
			} `json:"artifact"`
			DockerImage string `json:"docker-image"`
			Ports       struct {
				InternalDocker int `json:"internal-docker"`
				ExternalDocker int `json:"external-docker"`
			} `json:"ports"`
			Builder struct {
				Name        string `json:"name"`
				Version     string `json:"version"`
				DockerImage string `json:"docker-image"`
			} `json:"builder"`
		} `json:"application"`

		Database struct {
			Name         string `json:"name"`
			Driver       string `json:"driver"`
			User         string `json:"user"`
			Password     string `json:"password"`
			RootPassword string `json:"root_password"`
			Ports        struct {
				InternalDocker int `json:"internal-docker"`
				ExternalDocker int `json:"external-docker"`
			} `json:"ports"`
		} `json:"database"`

		ExtraCommands struct {
			LocalScriptFolderDir string `json:"local-script-folder-dir"`
			Commands             struct {
				Before []struct {
					Command  string `json:"command"`
					Priority int    `json:"priority"`
				} `json:"before"`
				After []struct {
					Command  string `json:"command"`
					RunOrder int    `json:"run-order"`
				} `json:"after"`
			} `json:"commands"`
		} `json:"extra-commands"`

		CodeCheckout struct {
			Git struct {
				Enabled      bool   `json:"enabled"`
				RepoURL      string `json:"repo-url"`
				RepoUser     string `json:"repo-user"`
				RepoPassword string `json:"repo-password"`
			} `json:"git"`

			SCP struct {
				Enabled        bool `json:"enabled"`
				Configurations struct {
					Folder string `json:"folder"`
				} `json:"configurations"`
			} `json:"scp"`
		} `json:"code-checkout"`
	} `json:"artifact"`
}
