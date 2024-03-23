package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		Log("INFO", ".env file not found")
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		Log("FATAL", fmt.Sprintf("failed to connect database %s", err.Error()))
		os.Exit(1)
	}
	db.AutoMigrate(&UserModel{}, &DocumentModel{}, &DocumentCommentModel{})

	Log("INFO", "Running server on :7777")
	err = http.ListenAndServe(":7777", DocsManagerServer{db})
	if err != nil {
		Log("FATAL", fmt.Sprintf("%s", err.Error()))
		os.Exit(1)
	}

}
