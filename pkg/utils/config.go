package utils

// Get config path for local or docker
func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	} else if configPath == "local" {
		return "./config/config-local"
	}
	return "./config"
}
