package main

import (
	"net/http"

	"gorm.io/gorm"
)

type DocsManagerServer struct {
	db *gorm.DB
}

func (svr DocsManagerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.URL.Path {
	case "/":
		w.WriteHeader(http.StatusOK)
		return

	case "/signup":
		svr.RoutePostSignUp(w, r)
		return

	case "/login":
		svr.RoutePostLogin(w, r)
		return

	case "/upload":
		svr.RoutePostUploadDocument(w, r)
		return

	case "/docs":
		svr.RouteGetDocuments(w, r)
		return

	default:
		w.WriteHeader(http.StatusNotFound)
		return

	}

}
