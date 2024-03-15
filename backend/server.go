package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gorm.io/gorm"
)

type DocsManagerServer struct {
	db *gorm.DB
}

func (svr DocsManagerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.URL.Path {
	case "/signup":
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			Log("ERROR", fmt.Sprintf("error parsing signup payload %s", err.Error()))

		}
		user := &UserModel{}
		json.Unmarshal(payload, user)

		u, _ := json.MarshalIndent(user, "", "  ")
		fmt.Println(fmt.Sprintf("%s", u))

		// svr.db.Create()
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "login")
	case "/":
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "ok")
		return
	case "/docs":
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "docs")
		return

	case "/login":
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "login")
	default:
		w.WriteHeader(http.StatusNotFound)
		return

	}

}
