package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gorm.io/gorm"
)

type DocsManagerServer struct {
	db  *gorm.DB
	mux *http.ServeMux
}

func (svr DocsManagerServer) SignUpRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		payload, err := io.ReadAll(r.Body)
		if err != nil {
			Log("ERROR", fmt.Sprintf("RoutePostSignUp - error parsing signup payload %s", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := signUpUser(payload, svr.db)
		if err != nil {
			Log("ERROR", fmt.Sprintf("RoutePostSignUp - error signing up user %s", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "plain/text")
		fmt.Fprintf(w, user.JWT)
	}
}

func (svr DocsManagerServer) LoginRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		payload, err := io.ReadAll(r.Body)
		if err != nil {
			Log("ERROR", fmt.Sprintf("RoutePostLogin - error parsing signup payload %s", err.Error()))
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
		w.Header().Set("Content-Type", "plain/text")
		fmt.Fprintf(w, user.JWT)
	}
}

func (svr DocsManagerServer) UploadRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		err := r.ParseMultipartForm(10 * 1024 * 1024) // 10 MB limit
		if err != nil {
			Log("ERROR", fmt.Sprintf("RoutePostUploadDocument - %s", err.Error()))
			http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()), http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			Log("ERROR", fmt.Sprintf("RoutePostUploadDocument - %s", err.Error()))
			http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()), http.StatusBadRequest)
			return
		}
		defer file.Close()

		doc := &DocumentModel{}
		uIdCtx := r.Context().Value("userId")
		uid, ok := uIdCtx.(uint)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		doc.UserID = uid
		err = uploadDocument(doc, file, handler, svr.db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()))
			return
		}

		res, _ := doc.toJson()
		fmt.Fprintf(w, res)
	}
}

func (svr DocsManagerServer) DocumentsRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		uIdCtx := r.Context().Value("userId")
		uid, ok := uIdCtx.(uint)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		docs := []DocumentModel{}
		svr.db.Find(&docs, "user_id = ?", uid)

		j, _ := json.Marshal(docs)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(j))
	}
}

func (svr DocsManagerServer) DocumentRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		qs := r.URL.Query()
		docUUID := qs["uuid"]
		if docUUID == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		http.ServeFile(w, r, fmt.Sprintf("docs/%s.pdf", docUUID[0]))
	}
}

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
