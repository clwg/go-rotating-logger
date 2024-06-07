# Golang Rotating Logger
This library provides a rotating logger in Go. It creates log files that rotate based on time and line count, and compresses the old log file on rotation.


# Usage

Configure the logger using the LoggerConfig struct:

```go
import (
	rotatinglogger "github.com/clwg/go-rotating-logger"
)

config.LoggerConfig = rotatinglogger.LoggerConfig{
    FilenamePrefix: *filenamePrefix,
    LogDir:         *logDir,
    MaxLines:       *maxLines,
    RotationTime:   time.Duration(*rotationTime) * time.Minute,
}
```

LogDir: The directory where the log files will be stored.
FilenamePrefix: The prefix for the log file names.
MaxLines: The maximum number of lines that a log file can have before it's rotated.
RotationTime: The duration after which the log file will be rotated.

Initialize the logger

```go
logger, err := rotatinglogger.NewLogger(config)
if err != nil {
    log.Fatalf("could not create logger: %v", err)
}
```

Log a message

```go
err := logger.Log(myObject)
if err != nil {
    log.Fatalf("could not log object: %v", err)
}
```