package log

import (
	"os"
	"testing"
)

func TestInfof(t *testing.T) {
	setupLogger(true, false, nil)
	defer Close()

	Infof("This is an %s message", "info")
	if len(logInstance.logStrs) == 0 || logInstance.logStrs[0] != "INFO : This is an info message\n" {
		t.Errorf("Expected INFO log message, got: %v", logInstance.logStrs)
	}
}

func TestWarnf(t *testing.T) {
	setupLogger(true, false, nil)
	defer Close()

	Warnf("This is a %s message", "warning")
	if len(logInstance.logStrs) == 0 || logInstance.logStrs[0] != "WARN : This is a warning message\n" {
		t.Errorf("Expected WARN log message, got: %v", logInstance.logStrs)
	}
}

func TestErrorf(t *testing.T) {
	setupLogger(true, false, nil)
	defer Close()

	Errorf("This is an %s message", "error")
	if len(logInstance.logStrs) == 0 || logInstance.logStrs[0] != "ERROR: This is an error message\n" {
		t.Errorf("Expected ERROR log message, got: %v", logInstance.logStrs)
	}
}

func TestDebugf(t *testing.T) {
	setupLogger(true, false, nil)
	defer Close()

	level = DEBUG
	Debugf("This is a %s message", "debug")
	if len(logInstance.logStrs) == 0 || logInstance.logStrs[0] != "DEBUG: This is a debug message\n" {
		t.Errorf("Expected DEBUG log message, got: %v", logInstance.logStrs)
	}
}

func TestFatalf(t *testing.T) {
	setupLogger(true, false, nil)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on Fatalf, but did not panic")
		}
		Close()
	}()

	Fatalf("This is a %s message", "fatal")
}

func TestClose(t *testing.T) {
	file, err := os.CreateTemp("", "log_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	setupLogger(false, false, file)
	Close()

	if logInstance != nil {
		t.Errorf("Expected logInstance to be nil after Close, got: %v", logInstance)
	}
}

func setupLogger(inMemory, stdOut bool, logFile *os.File) {
	logInstance = &logger{
		inMemory:           inMemory,
		stdOut:             stdOut,
		logFile:            logFile,
		formatLogMessageFn: formatLogMessageWithoutTime,
	}
	isInitialized = true
}
