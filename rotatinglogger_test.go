package rotatinglogger_test

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	rotatinglogger "github.com/clwg/go-rotating-logger"
)

const (
	testLogDir         = "./test_logs"
	testFilenamePrefix = "test_log"
)

func setupTestLogger(maxLines int, rotationTime time.Duration) (*rotatinglogger.Logger, error) {
	config := rotatinglogger.LoggerConfig{
		LogDir:         testLogDir,
		LogFormat:      rotatinglogger.FormatJSON,
		FilenamePrefix: testFilenamePrefix,
		MaxLines:       maxLines,
		RotationTime:   rotationTime,
	}
	return rotatinglogger.NewLogger(config)
}

func teardownTestLogger() {
	os.RemoveAll(testLogDir)
}

func TestLogger_SizeBasedRotation(t *testing.T) {
	teardownTestLogger()
	defer teardownTestLogger()

	logger, err := setupTestLogger(2, 24*time.Hour) // Set max lines to 2 for test
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	for i := 0; i < 5; i++ { // Log more lines than maxLines
		err := logger.Log(map[string]string{"message": fmt.Sprintf("log entry %d", i)})
		if err != nil {
			t.Fatalf("Failed to log entry: %v", err)
		}
		time.Sleep(1000 * time.Millisecond) // Slight delay to ensure different timestamp in filename
	}

	archiveDir := filepath.Join(testLogDir, "archive")
	files, err := os.ReadDir(archiveDir)
	if err != nil {
		t.Fatalf("Failed to read archive directory: %v", err)
	}

	expectedFiles := 2 // Based on 5 log entries with maxLines of 2
	if len(files) != expectedFiles {
		t.Fatalf("Expected %d archived log files, got %d", expectedFiles, len(files))
	}
}

func TestLogger_TimeBasedRotation(t *testing.T) {
	teardownTestLogger()
	defer teardownTestLogger()

	logger, err := setupTestLogger(100, 1*time.Second) // Set rotation time to 1 second for test
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	err = logger.Log(map[string]string{"message": "log entry 1"})
	if err != nil {
		t.Fatalf("Failed to log entry: %v", err)
	}

	time.Sleep(2 * time.Second) // Wait for rotation time to pass

	err = logger.Log(map[string]string{"message": "log entry 2"})
	if err != nil {
		t.Fatalf("Failed to log entry: %v", err)
	}

	archiveDir := filepath.Join(testLogDir, "archive")
	files, err := os.ReadDir(archiveDir)
	if err != nil {
		t.Fatalf("Failed to read archive directory: %v", err)
	}

	expectedFiles := 1
	if len(files) != expectedFiles {
		t.Fatalf("Expected %d archived log files, got %d", expectedFiles, len(files))
	}
}

func TestLogger_ArchiveContent(t *testing.T) {
	teardownTestLogger()
	defer teardownTestLogger()

	logger, err := setupTestLogger(2, 24*time.Hour) // Set max lines to 2 for test
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	for i := 0; i < 3; i++ {
		err := logger.Log(map[string]string{"message": fmt.Sprintf("log entry %d", i)})
		if err != nil {
			t.Fatalf("Failed to log entry: %v", err)
		}
	}

	time.Sleep(1000 * time.Millisecond) // Slight delay to ensure different timestamp in filename

	archiveDir := filepath.Join(testLogDir, "archive")
	files, err := os.ReadDir(archiveDir)
	if err != nil {
		t.Fatalf("Failed to read archive directory: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("No archived log files found")
	}

	archiveFilePath := filepath.Join(archiveDir, files[0].Name())
	archiveFile, err := os.Open(archiveFilePath)
	if err != nil {
		t.Fatalf("Failed to open archived log file: %v", err)
	}
	defer archiveFile.Close()

	gzipReader, err := gzip.NewReader(archiveFile)
	if err != nil {
		t.Fatalf("Failed to create gzip reader: %v", err)
	}
	defer gzipReader.Close()

	_, err = io.Copy(io.Discard, gzipReader)
	if err != nil {
		t.Fatalf("Failed to read archived log content: %v", err)
	}
}
