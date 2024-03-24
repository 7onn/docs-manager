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

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{}) // todo: use stateful db
	if err != nil {
		Log("FATAL", fmt.Sprintf("failed to connect database %s", err.Error()))
		os.Exit(1)
	}
	db.AutoMigrate(&UserModel{}, &DocumentModel{}, &DocumentCommentModel{})

	svr := &DocsManagerServer{
		db:  db,
		mux: http.NewServeMux(),
	}

	svr.mux.HandleFunc("/", svr.HomeRoute())
	svr.mux.HandleFunc("/healthz", svr.HealthCheckRoute())
	svr.mux.HandleFunc("/signup", svr.SignUpRoute())
	svr.mux.HandleFunc("/login", svr.LoginRoute())
	svr.mux.HandleFunc("/upload", verifyJWT(svr.UploadRoute()))
	svr.mux.HandleFunc("/docs", verifyJWT(svr.DocumentsRoute()))
	svr.mux.HandleFunc("/doc", svr.DocumentRoute())

	p := os.Getenv("SERVER_PORT")
	if p == "" {
		p = "7777"
	}
	Log("INFO", fmt.Sprintf("Running server on :%s", p))
	err = http.ListenAndServe(fmt.Sprintf(":%s", p), withCORS(svr.mux))
	if err != nil {
		Log("FATAL", fmt.Sprintf("%s", err.Error()))
		os.Exit(1)
	}
}
