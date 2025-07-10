package log

import (
	"fmt"
	"os"
	"time"
)

var (
	logInstance   *logger
	isInitialized bool = false
)

type logger struct {
	logFile  *os.File
	inMemory bool
	stdOut   bool
	logStrs  []string
}

func (l *logger) init() error {
	return nil
}

func printLog(typeStr string, format string, args ...any) {
	logStr := fmt.Sprintf("%s [%s]: ", typeStr, time.Now().Format("2006-01-02 15:04:05.000")) + fmt.Sprintf(format, args...)

	// add newline at the end if not already present
	if len(logStr) > 0 && logStr[len(logStr)-1] != '\n' {
		logStr += "\n"
	}

	if logInstance.inMemory {
		logInstance.logStrs = append(logInstance.logStrs, logStr)
	} else if logInstance.stdOut {
		fmt.Print(logStr)
	}
	if logInstance.logFile != nil {
		// Ignore error for simplicity
		logInstance.logFile.WriteString(logStr)
	}
}

func Infof(format string, args ...any) {
	if !isInitialized {
		return
	}

	printLog("Info ", format, args...)
}

func Warnf(format string, args ...any) {
	if !isInitialized {
		return
	}

	printLog("Warn ", format, args...)
}

func Errorf(format string, args ...any) {
	if !isInitialized {
		return
	}

	printLog("Error", format, args...)
}

func Debugf(format string, args ...any) {
	if !isInitialized {
		return
	}

	printLog("Debug", format, args...)
}

func Fatalf(format string, args ...any) {
	if !isInitialized {
		return
	}

	printLog("Fatal", format, args...)

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
