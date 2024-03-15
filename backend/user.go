package main

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Id        string `json:"id",gorm:"primaryKey"`
	Name      string `json:"name"`
	Email     string `json:"Email"`
	Validated bool   `json:omit`
}

func (u *UserModel) BeforeCreate(tx *gorm.DB) error {
	u.Id = uuid.NewString()
	return nil
}

func signUpUser(user UserModel) {

}
