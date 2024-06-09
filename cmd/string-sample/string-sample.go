package main

import (
	"fmt"
	"time"

	rotatinglogger "github.com/clwg/go-rotating-logger"
)

type SampleData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	config := rotatinglogger.LoggerConfig{
		LogDir:         "./logs",
		FilenamePrefix: "sample",
		MaxLines:       100,
		LogFormat:      rotatinglogger.FormatText,
		RotationTime:   time.Duration(60) * time.Minute,
	}

	logger, err := rotatinglogger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}

	counter := 0
	for i := 0; i < 726; i++ {
		fmt.Printf("Logging data %d\n", i)
		data := fmt.Sprintf("Sample Name %d", i)

		err := logger.Log(data)
		if err != nil {
			fmt.Printf("Failed to log data: %v\n", err)
			return
		}

		counter++
		// Sleep to avoid trying to write files with the same timestamp
		if counter == 100 {
			time.Sleep(1 * time.Second)
			counter = 0
		}
	}
}
