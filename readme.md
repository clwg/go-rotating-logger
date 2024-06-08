# Golang Rotating Logger
This library provides a rotating logger in Go. It creates log files that rotate based on time and line count, and compresses the old log file on rotation.


## Usage

import the library

```
import (
	rotatinglogger "github.com/clwg/go-rotating-logger"
)
```

Initialize the logger setting the file and path as well as time and line based rotation parameters.

```go
config := rotatinglogger.LoggerConfig{
    LogDir:         "./logs",
    FilenamePrefix: "sample",
    MaxLines:       100,
    RotationTime:   time.Duration(60) * time.Minute,
}

logger, err := rotatinglogger.NewLogger(config)
if err != nil {
    fmt.Printf("Failed to create logger: %v\n", err)
    return
}
```

Then just call the logger to log messages

```go
err := logger.Log(data)
if err != nil {
    fmt.Printf("Failed to log data: %v\n", err)
    return
}
```

