package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (svr DocsManagerServer) RoutePostSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		Log("ERROR", fmt.Sprintf("error parsing signup payload %s", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := signUpUser(payload, svr.db)
	if err != nil {
		Log("ERROR", fmt.Sprintf("error signing up user %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(user)
	fmt.Fprintf(w, fmt.Sprintf("%s", j))
}

func (svr DocsManagerServer) RoutePostLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		Log("ERROR", fmt.Sprintf("error parsing signup payload %s", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()))
		return
	}

	user, err := loginUser(payload, svr.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	j, _ := user.toJson()
	fmt.Fprintf(w, j)
}

func (svr DocsManagerServer) RouteGetDocuments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	docs := []DocumentModel{}
	svr.db.Find(&docs)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(docs)
	fmt.Fprintf(w, fmt.Sprintf("%s", j))
}

func (svr DocsManagerServer) RoutePostUploadDocument(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseMultipartForm(10 * 1024 * 1024) // 10 MB limit
	if err != nil {
		Log("ERROR", err.Error())
		http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		Log("ERROR", err.Error())
		http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()), http.StatusBadRequest)
		return
	}
	defer file.Close()

	doc := &DocumentModel{}
	doc.UserID = 1 // todo: get user from jwt cookie

	err = uploadDocument(doc, file, handler, svr.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()))
		return
	}

	res, _ := doc.toJson()
	fmt.Fprintf(w, res)
}
