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

	p := os.Getenv("SERVER_PORT")
	if p == "" {
		p = "7777"
	}

	svr := &DocsManagerServer{
		db:  db,
		mux: http.NewServeMux(),
	}

	svr.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})

	svr.mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		u := UserModel{}
		test := db.First(&u)
		if test.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	})

	svr.mux.HandleFunc("/signup", svr.SignUpRoute())
	svr.mux.HandleFunc("/login", svr.LoginRoute())
	svr.mux.HandleFunc("/upload", verifyJWT(svr.UploadRoute()))
	svr.mux.HandleFunc("/docs", verifyJWT(svr.DocumentsRoute()))
	svr.mux.HandleFunc("/doc", svr.DocumentRoute())

	Log("INFO", fmt.Sprintf("Running server on :%s", p))
	err = http.ListenAndServe(fmt.Sprintf(":%s", p), withCORS(svr.mux))
	if err != nil {
		Log("FATAL", fmt.Sprintf("%s", err.Error()))
		os.Exit(1)
	}
}
