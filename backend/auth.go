package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type UserJWTClaims struct {
	UUID   string `json:"uuid"`
	UserID uint   `json:"userId"`
	jwt.StandardClaims
}

func (u *UserModel) setJWT() (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, UserJWTClaims{
		UUID:   u.UUID.String(),
		UserID: u.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev"
	}
	s, err := t.SignedString([]byte(jwtSecret))
	if err != nil {
		Log("ERROR", fmt.Sprintf("setJWT - %s", err.Error()))
		return "", err
	}

	u.JWT = s
	return s, nil
}

func signUpUser(payload []byte, db *gorm.DB) (*UserModel, error) {
	u := &UserModel{}
	json.Unmarshal(payload, u)
	db.Create(u)
	_, err := u.setJWT()
	if err != nil {
		return u, err
	}
	db.Updates(u)
	return u, nil
}

func loginUser(payload []byte, db *gorm.DB) (*UserModel, error) {
	u := &UserModel{}
	u.fromJson(payload)
	pwd := u.Password

	db.First(u, "email = ?", u.Email)
	if u.UUID.String() == "" {
		return nil, errors.New("Unexistent user")
	}

	if !u.CheckPasswordHash(pwd) {
		return nil, errors.New("Wrong password")
	}
	_, err := u.setJWT()
	if err != nil {
		return u, err
	}
	db.Updates(u)

	return u, nil
}

func verifyJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header["Authorization"]
		if auth == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		jwtToken := auth[0]
		jwtToken, err := url.QueryUnescape(jwtToken)
		if err != nil {
			return
		}
		claims := &UserJWTClaims{}
		token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
			jwtSecret := os.Getenv("JWT_SECRET")
			if jwtSecret == "" {
				jwtSecret = "dev"
			}
			return []byte(jwtSecret), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		claims, ok := token.Claims.(*UserJWTClaims)
		if ok {
			ctx := context.WithValue(r.Context(), "uuid", claims.UUID)
			ctx = context.WithValue(r.Context(), "userId", claims.UserID)
			if !token.Valid {
				cookie := http.Cookie{}
				cookie.Name = "jwt"
				cookie.Value = ""
				cookie.Secure = false
				cookie.HttpOnly = true
				cookie.Path = "/"
				http.SetCookie(w, &cookie)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next(w, r.WithContext(ctx))
		}
	}
}
