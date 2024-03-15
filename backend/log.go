package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type LogData struct {
	Timestamp time.Time `json:"timestamp"`
	Severity  string    `json:"severity"`
	Message   string    `json:"message"`
}

func Log(severity, message string) {
	if severity == "DEBUG" && os.Getenv("LOG_DEBUG") != "true" {
		return
	}
	logData := LogData{
		Timestamp: time.Now(),
		Severity:  severity,
		Message:   message,
	}

	jsonData, err := json.Marshal(logData)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}
