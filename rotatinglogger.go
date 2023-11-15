package rotatinglogger

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LoggerConfig struct {
	LogDir         string
	FilenamePrefix string
	MaxLines       int
	RotationTime   time.Duration
}

type Logger struct {
	config      LoggerConfig
	currentFile *os.File
	lineCount   int
	mu          sync.Mutex
	ticker      *time.Ticker
}

func NewLogger(config LoggerConfig) (*Logger, error) {
	logger := &Logger{
		config: config,
		ticker: time.NewTicker(config.RotationTime),
	}
	err := logger.rotateFile()
	if err != nil {
		return nil, err
	}
	go logger.timeBasedRotation()
	return logger, nil
}

func (l *Logger) Log(obj interface{}) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	line, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(l.currentFile, string(line))
	if err != nil {
		return err
	}

	l.lineCount++
	if l.lineCount >= l.config.MaxLines {
		return l.rotateFileAndResetTimer()
	}

	return nil
}

func (l *Logger) rotateFile() error {
	// Ensure the log directory exists
	if _, err := os.Stat(l.config.LogDir); os.IsNotExist(err) {
		if err := os.MkdirAll(l.config.LogDir, 0755); err != nil {
			return err
		}
	}

	if l.currentFile != nil {
		archiveDir := filepath.Join(l.config.LogDir, "archive")
		if _, err := os.Stat(archiveDir); os.IsNotExist(err) {
			if err := os.MkdirAll(archiveDir, 0755); err != nil {
				return err
			}
		}
		gzipFilename := filepath.Join(archiveDir, fmt.Sprintf("%s.gz", filepath.Base(l.currentFile.Name())))
		gzipFile, err := os.Create(gzipFilename)
		if err != nil {
			return err
		}
		defer gzipFile.Close()

		writer := gzip.NewWriter(gzipFile)
		defer writer.Close()

		l.currentFile.Seek(0, 0)
		if _, err := io.Copy(writer, l.currentFile); err != nil {
			return err
		}
		writer.Close()
		gzipFile.Close()

		originalFileName := l.currentFile.Name()
		l.currentFile.Close()

		if err := os.Remove(originalFileName); err != nil {
			return err
		}
	}

	filename := fmt.Sprintf("%s_%s.log", l.config.FilenamePrefix, time.Now().Format(time.RFC3339))
	fullPath := filepath.Join(l.config.LogDir, filename)
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	l.currentFile = file
	l.lineCount = 0
	return nil
}

func (l *Logger) rotateFileAndResetTimer() error {
	err := l.rotateFile()
	if err != nil {
		return err
	}
	l.ticker.Reset(l.config.RotationTime)
	return nil
}

func (l *Logger) timeBasedRotation() {
	for range l.ticker.C {
		l.mu.Lock()
		l.rotateFileAndResetTimer()
		l.mu.Unlock()
	}
}
