package log

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	logInstance   *logger
	isInitialized bool = false
	level              = INFO // Default log level
)

const (
	FATAL = iota
	ERROR
	WARN
	INFO
	DEBUG
)

type fmtFunc func(string, string, ...any) string

type logger struct {
	mx                 sync.Mutex
	logFile            *os.File
	inMemory           bool
	stdOut             bool
	logStrs            []string
	formatLogMessageFn fmtFunc
}

func (l *logger) init() error {
	return nil
}

func formatLogMessageWithTime(typeStr string, format string, args ...any) string {
	return fmt.Sprintf("%s [%s]: ", typeStr, time.Now().Format("2006-01-02 15:04:05.000")) + fmt.Sprintf(format, args...)
}

func formatLogMessageWithoutTime(typeStr string, format string, args ...any) string {
	return typeStr + ": " + fmt.Sprintf(format, args...)
}

func printLog(typeStr string, format string, args ...any) {
	logStr := logInstance.formatLogMessageFn(typeStr, format, args...)

	// add newline at the end if not already present
	if len(logStr) > 0 && logStr[len(logStr)-1] != '\n' {
		logStr += "\n"
	}

	if logInstance.inMemory {
		logInstance.mx.Lock()
		defer logInstance.mx.Unlock()

		logInstance.logStrs = append(logInstance.logStrs, logStr)
	} else if logInstance.stdOut {
		fmt.Print(logStr)
	}
	if logInstance.logFile != nil {
		logInstance.mx.Lock()
		defer logInstance.mx.Unlock()

		// Ignore error for simplicity
		logInstance.logFile.WriteString(logStr)
	}
}

func Infof(format string, args ...any) {
	if !isInitialized || level < INFO {
		return
	}

	printLog("INFO ", format, args...)
}

func Warnf(format string, args ...any) {
	if !isInitialized || level < WARN {
		return
	}

	printLog("WARN ", format, args...)
}

func Errorf(format string, args ...any) {
	if !isInitialized || level < ERROR {
		return
	}

	printLog("ERROR", format, args...)
}

func Debugf(format string, args ...any) {
	if !isInitialized || level < DEBUG {
		return
	}

	printLog("DEBUG", format, args...)
}

func Fatalf(format string, args ...any) {
	if !isInitialized || level < FATAL {
		return
	}

	printLog("FATAL", format, args...)

	// Exit the program after logging fatal error
	panic("Fatal error occurred, exiting program")
}

func Close() {
	if !isInitialized {
		return
	}
	if logInstance.logFile != nil {
		err := logInstance.logFile.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close log file: %v\n", err)
		}
		logInstance.logFile = nil
	}

	if logInstance.inMemory && logInstance.stdOut {
		for _, logStr := range logInstance.logStrs {
			fmt.Print(logStr)
		}
	}
	isInitialized = false
	logInstance = nil
}
