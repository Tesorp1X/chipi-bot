package config

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/Tesorp1X/chipi-bot/utils"
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
	}

	admins, err := extractAdmins()
	if err != nil {
		return nil, fmt.Errorf(
			"error in InitConfig(): couldn't extract admins: %v",
			err,
		)
	}

	dbPath, err := extractDbPath()
	if err != nil {
		return nil, fmt.Errorf(
			"error in InitConfig(): couldn't extract db path: %v",
			err,
		)
	}

	downloadPath, err := extractDownloadPath()
	if err != nil {
		return nil, fmt.Errorf(
			"error in InitConfig(): couldn't extract download path: %v",
			err,
		)
	}

	logPath, err := extractLogPath()
	if err != nil {
		return nil, fmt.Errorf(
			"error in InitConfig(): couldn't extract log path: %v",
			err,
		)
	}

	return &Config{
		ApiKey:       apiKey,
		Admins:       admins,
		VerboseDebug: verboseDebug,
		DbPath:       dbPath,
		DownloadPath: downloadPath,
		LogPath:      logPath,
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

func extractAdmins() ([]int64, error) {
	var adminsList []int64

	adminsListStr := os.Getenv("ADMIN_LIST")
	if len(adminsListStr) == 0 {
		return nil, fmt.Errorf(
			"error: error in extractAdmins(): empty ADMIN_LIST",
		)
	}

	adminsList, err := utils.ExtractAdminsIDs(adminsListStr)
	if err != nil {
		return nil, fmt.Errorf(
			"error in extractAdmins(): couldn't extract admins ids from '%s': %v",
			adminsListStr,
			err,
		)
	}

	return adminsList, nil
}

func extractDbPath() (string, error) {
	dbPath := os.Getenv("DB_PATH")
	if !fs.ValidPath(dbPath) || !strings.HasSuffix(dbPath, ".db") {
		return "", fmt.Errorf(
			"error in extractDbPath(): db path '%s' is invalid",
			dbPath,
		)
	}

	return dbPath, nil
}

func extractDownloadPath() (string, error) {
	dbPath := os.Getenv("DOWNLOAD_PATH")
	// TODO: also check that path is a dir and not a file
	if !fs.ValidPath(dbPath) {
		return "", fmt.Errorf(
			"error in extractDownloadPath(): download path '%s' is invalid",
			dbPath,
		)
	}

	return dbPath, nil
}

func extractLogPath() (string, error) {
	dbPath := os.Getenv("LOG_PATH")
	// TODO: also check that path is a dir and not a file
	if !fs.ValidPath(dbPath) {
		return "", fmt.Errorf(
			"error in extractLogPath(): log path '%s' is invalid",
			dbPath,
		)
	}

	return dbPath, nil
}
