package config

import (
	"os"
	"path"
	"strconv"
)

var (
	baseDir string = ".gotchat"

	// Windows specific
	windowsBaseDirPtr *string
)

func GetBaseDir() string {
	if baseDirEnv, ok := os.LookupEnv("GOTCHAT_BASE_DIR"); ok {
		return baseDirEnv
	}

	if windowsBaseDirPtr != nil {
		return *windowsBaseDirPtr
	}

	return baseDir
}

func GetRootDir() string {
	if rootDir, ok := os.LookupEnv("GOTCHAT_ROOT_DIR"); ok {
		return rootDir
	}

	gotchatDir := GetBaseDir()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return gotchatDir
	}

	return path.Join(homeDir, gotchatDir)
}

func GetDataStorageFilePath() string {
	var dataStorageFileName string
	var ok bool

	if dataStorageFileName, ok = os.LookupEnv("GOTCHAT_DATA_STORAGE_FILE_NAME"); !ok {
		dataStorageFileName = "chat.db"
	}

	return path.Join(GetRootDir(), dataStorageFileName)
}

func GetServerPort() int {
	port := 7665 // default port

	if portStr, ok := os.LookupEnv("GOTCHAT_SERVER_PORT"); ok {
		var err error
		if port, err = strconv.Atoi(portStr); err != nil {
			port = 7665 // default port
		}
	}

	return port
}
