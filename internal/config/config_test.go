package config

import (
	"os"
	"path"
	"testing"
)

func TestGetBaseDir(t *testing.T) {
	originalEnv := os.Getenv("GOTCHAT_BASE_DIR")
	defer os.Setenv("GOTCHAT_BASE_DIR", originalEnv)

	// Test with environment variable set
	os.Setenv("GOTCHAT_BASE_DIR", "/custom/base/dir")
	if got := GetBaseDir(); got != "/custom/base/dir" {
		t.Errorf("GetBaseDir() = %v, want %v", got, "/custom/base/dir")
	}

	// Test without environment variable
	os.Unsetenv("GOTCHAT_BASE_DIR")
	windowsBaseDirPtr = nil
	if got := GetBaseDir(); got != ".gotchat" {
		t.Errorf("GetBaseDir() = %v, want %v", got, ".gotchat")
	}
}

func TestGetRootDir(t *testing.T) {
	originalEnv := os.Getenv("GOTCHAT_ROOT_DIR")
	defer os.Setenv("GOTCHAT_ROOT_DIR", originalEnv)

	// Test with environment variable set
	os.Setenv("GOTCHAT_ROOT_DIR", "/custom/root/dir")
	if got := GetRootDir(); got != "/custom/root/dir" {
		t.Errorf("GetRootDir() = %v, want %v", got, "/custom/root/dir")
	}

	// Test without environment variable
	os.Unsetenv("GOTCHAT_ROOT_DIR")
	os.Unsetenv("GOTCHAT_BASE_DIR")
	homeDir, _ := os.UserHomeDir()
	expected := path.Join(homeDir, ".gotchat")
	if got := GetRootDir(); got != expected {
		t.Errorf("GetRootDir() = %v, want %v", got, expected)
	}
}

func TestGetDataStorageFilePath(t *testing.T) {
	originalEnv := os.Getenv("GOTCHAT_DATA_STORAGE_FILE_NAME")
	defer os.Setenv("GOTCHAT_DATA_STORAGE_FILE_NAME", originalEnv)

	// Test with environment variable set
	os.Setenv("GOTCHAT_DATA_STORAGE_FILE_NAME", "custom.db")
	expected := path.Join(GetRootDir(), "custom.db")
	if got := GetDataStorageFilePath(); got != expected {
		t.Errorf("GetDataStorageFilePath() = %v, want %v", got, expected)
	}

	// Test without environment variable
	os.Unsetenv("GOTCHAT_DATA_STORAGE_FILE_NAME")
	expected = path.Join(GetRootDir(), "chat.db")
	if got := GetDataStorageFilePath(); got != expected {
		t.Errorf("GetDataStorageFilePath() = %v, want %v", got, expected)
	}
}

func TestGetServerPort(t *testing.T) {
	originalEnv := os.Getenv("GOTCHAT_SERVER_PORT")
	defer os.Setenv("GOTCHAT_SERVER_PORT", originalEnv)

	// Test with environment variable set to a valid port
	os.Setenv("GOTCHAT_SERVER_PORT", "8080")
	if got := GetServerPort(); got != 8080 {
		t.Errorf("GetServerPort() = %v, want %v", got, 8080)
	}

	// Test with environment variable set to an invalid port
	os.Setenv("GOTCHAT_SERVER_PORT", "invalid")
	if got := GetServerPort(); got != 7665 {
		t.Errorf("GetServerPort() = %v, want %v", got, 7665)
	}

	// Test without environment variable
	os.Unsetenv("GOTCHAT_SERVER_PORT")
	if got := GetServerPort(); got != 7665 {
		t.Errorf("GetServerPort() = %v, want %v", got, 7665)
	}
}
