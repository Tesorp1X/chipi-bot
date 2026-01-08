package config

import (
	"fmt"
	"github.com/joho/godotenv"
)

type Config struct {
	// "API_KEY" in env
	ApiKey string
	// "ADMIN_LIST" in env
	Admins []int64
	// as flag
	VerboseDebug bool
	// "DB_PATH" in env
	DbPath string
	// "DOWNLOAD_PATH" in env
	DownloadPath string
	// "LOG_PATH" in env
	LogPath string
}

func InitConfig(verboseDebug bool) (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		panic(
			fmt.Errorf(
				"error in InitConfig(): couldn't load a '.env' file: %v",
				err,
			))
	}

	apiKey, err := extractApiKey()
	if err != nil {
		return nil, fmt.Errorf(
			"error in InitConfig(): couldn't extract apiKey: %v",
			err,
		)
	return &Config{
		ApiKey:       apiKey,
	}, nil
}

// Extracts tg-bot api key from the environment variable API_KEY.
// Returns a non-nil error in case, if there were no such variable.
func extractApiKey() (string, error) {
	apiKey := os.Getenv("API_KEY")
	if len(apiKey) == 0 {
		return "", fmt.Errorf(
			"error in extractApiKey(): empty API_KEY.",
		)
	}
	return apiKey, nil
}
	}
