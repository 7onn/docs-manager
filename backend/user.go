package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	ID          uint            `json:"-" gorm:"primaryKey"`
	UUID        uuid.UUID       `json:"uuid" gorm:"unique" gorm:"index"`
	Name        string          `json:"name"`
	Email       string          `json:"email" gorm:"unique" gorm:"index"`
	Password    string          `json:"-"`
	JWT         string          `json:"-"`
	CreatedAt   time.Time       `json:"-"`
	UpdatedAt   time.Time       `json:"-"`
	ActivatedAt sql.NullTime    `json:"-"`
	Documents   []DocumentModel `gorm:"foreignKey:UserID"`
}

func (u *UserModel) BeforeCreate(tx *gorm.DB) error {
	u.UUID = uuid.New()
	return u.HashPassword()
}

func (u *UserModel) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *UserModel) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *UserModel) toJson() (string, error) {
	jb, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", jb), nil
}

func (u *UserModel) fromJson(data []byte) error {
	err := json.Unmarshal(data, u)
	if err != nil {
		return err
	}
	return nil
}
