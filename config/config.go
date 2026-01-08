package config

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
