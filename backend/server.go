package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type DocsManagerServer struct {
	db *gorm.DB
}

func (svr DocsManagerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // todo: allow actual hostname
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Println(r.URL.Path)
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
		if !validJWT(w, r) {
			return
		}
		svr.RoutePostUploadDocument(w, r)
		return

	case "/docs":
		if !validJWT(w, r) {
			return
		}
		svr.RouteGetDocuments(w, r)
		return

	case "/doc":
		// if !validJWT(w, r) {
		// 	return
		// }
		svr.RouteGetDocument(w, r)
		return

	default:
		w.WriteHeader(http.StatusNotFound)
		return

	}

}

func validJWT(w http.ResponseWriter, r *http.Request) bool {
	jwtToken := r.Header["Authorization"][0]
	jwtToken, err := url.QueryUnescape(jwtToken)
	if err != nil {
		return false
	}
	claims := &UserJWTClaims{}
	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return false
	}
	if claims, ok := token.Claims.(*UserJWTClaims); ok {
		fmt.Println(claims.UUID, claims.StandardClaims.ExpiresAt)
	}

	if !token.Valid {
		cookie := http.Cookie{}
		cookie.Name = "jwt"
		cookie.Value = ""
		cookie.Secure = false
		cookie.HttpOnly = true
		cookie.Path = "/"
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	return true
}

func (svr DocsManagerServer) RoutePostSignUp(w http.ResponseWriter, r *http.Request) {
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

func (svr DocsManagerServer) RoutePostLogin(w http.ResponseWriter, r *http.Request) {
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

func (svr DocsManagerServer) RouteGetDocument(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	qs := r.URL.Query()
	docUUID := qs["uuid"]
	fmt.Println(docUUID[0])
	if docUUID == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, fmt.Sprintf("docs/%s.pdf", docUUID[0]))
}

func (svr DocsManagerServer) RoutePostUploadDocument(w http.ResponseWriter, r *http.Request) {
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
