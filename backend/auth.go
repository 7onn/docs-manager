package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type UserJWTClaims struct {
	UUID string `json:"uuid"`
	jwt.StandardClaims
}

var jwtSecret = []byte("dev") //todo: use jwt secret from env

func (u *UserModel) setJWT() (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, UserJWTClaims{
		UUID: u.UUID.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	})

	s, err := t.SignedString(jwtSecret)
	if err != nil {
		Log("ERROR", fmt.Sprintf("setJWT - %s", err.Error()))
		return "", err
	}

	u.JWT = s
	fmt.Println(s)
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
