package log

import (
	"os"
	"testing"
)

func TestConfigure(t *testing.T) {
	builder := Configure()
	if builder == nil {
		t.Fatal("Expected LogBuilder instance, got nil")
	}
	if builder.level != INFO {
		t.Errorf("Expected default log level to be INFO, got %d", builder.level)
	}
}

func TestLogBuilder_Init(t *testing.T) {
	builder := Configure()
	err := builder.Init()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !isInitialized {
		t.Fatal("Expected logger to be initialized")
	}
}

func TestLogBuilder_InMemory(t *testing.T) {
	builder := Configure().InMemory()
	if !builder.logInstance.inMemory {
		t.Fatal("Expected inMemory to be true")
	}
}

func TestLogBuilder_StdOut(t *testing.T) {
	builder := Configure().StdOut()
	if !builder.logInstance.stdOut {
		t.Fatal("Expected stdOut to be true")
	}
}

func TestLogBuilder_File(t *testing.T) {
	filePath := "test.log"
	defer os.Remove(filePath)

	builder := Configure().File(filePath)
	if builder == nil {
		t.Fatal("Expected LogBuilder instance, got nil")
	}
	if builder.logInstance.logFile == nil {
		t.Fatal("Expected logFile to be initialized")
	}
}

func TestLogBuilder_Level(t *testing.T) {
	builder := Configure().Level(DEBUG)
	if builder == nil {
		t.Fatal("Expected LogBuilder instance, got nil")
	}
	if builder.level != DEBUG {
		t.Errorf("Expected log level to be DEBUG, got %d", builder.level)
	}

	invalidBuilder := Configure().Level(-1)
	if invalidBuilder != nil {
		t.Fatal("Expected nil for invalid log level, got LogBuilder instance")
	}
}

func TestLogBuilder_WithTimestamps(t *testing.T) {
	builder := Configure().WithTimestamps()
	if builder.logInstance.formatLogMessageFn == nil {
		t.Fatal("Expected formatLogMessageFn to be set")
	}
}

func TestLogBuilder_WithoutTimestamps(t *testing.T) {
	builder := Configure().WithoutTimestamps()
	if builder.logInstance.formatLogMessageFn == nil {
		t.Fatal("Expected formatLogMessageFn to be set")
	}
}
func TestLogBuilder_WithCustomFormat(t *testing.T) {
	// Test setting custom format function
	builder := Configure()

	customFormatter := func(level int, message string, params ...any) string {
		return "CUSTOM: " + message
	}

	result := builder.WithCustomFormat(customFormatter)

	if result == nil {
		t.Error("Expected WithCustomFormat to return LogBuilder, got nil")
	}

	if result != builder {
		t.Error("Expected WithCustomFormat to return the same builder instance")
	}
}

func TestLogBuilder_WithCustomFormat_NilFunction(t *testing.T) {
	// Test setting nil custom format function
	builder := Configure()

	result := builder.WithCustomFormat(nil)

	if result == nil {
		t.Error("Expected WithCustomFormat to return LogBuilder even with nil function, got nil")
	}

	if result != builder {
		t.Error("Expected WithCustomFormat to return the same builder instance")
	}

	if builder.logInstance.formatLogMessageFn == nil {
		t.Error("Expected formatLogMessageFn to remain unchanged when nil function is provided")
	}
}
