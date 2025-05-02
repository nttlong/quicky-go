package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Logger is an interface for logging
type Logger interface {
	Log(message string)
}

// FileLogger logs messages to a file
type FileLogger struct {
	logger *log.Logger
}

// NewFileLogger creates a new FileLogger that logs to logs/app.log
func NewFileLogger() (*FileLogger, error) {
	logDir := filepath.Join(".", "logs")
	logFile := filepath.Join(logDir, "app.log")

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Open or create log file
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Create logger
	logger := log.New(file, "ERP: ", log.LstdFlags)
	return &FileLogger{logger: logger}, nil
}

// Log writes a message to the log file
func (f *FileLogger) Log(message string) {
	f.logger.Println(message)
}
