package main

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentModel struct {
	gorm.Model
	ID       uint                   `json:"id"`
	UserID   uint                   `json:"userId"`
	User     UserModel              `json:"-" gorm:"foreignKey:UserID;references:ID"`
	UUID     uuid.UUID              `json:"uuid" gorm:"unique" gorm:"index"`
	Name     string                 `json:"name"`
	Comments []DocumentCommentModel `gorm:"foreignKey:DocumentID"`
}

type DocumentCommentModel struct {
	gorm.Model
	ID         uint
	DocumentID uint          `json:"docId"`
	UserID     uint          `json:"userId"`
	Document   DocumentModel `gorm:"foreignKey:DocumentID;references:ID"`
	User       UserModel     `gorm:"foreignKey:DocumentID;references:ID"`
	Message    string        `json:"message"`
	CreatedAt  time.Time     `json:"-"`
	UpdatedAt  time.Time     `json:"-"`
}

func (d *DocumentModel) BeforeCreate(tx *gorm.DB) error {
	d.UUID = uuid.New()
	return nil
}

func (d *DocumentModel) toJson() (string, error) {
	jb, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", jb), nil
}

func uploadDocument(doc *DocumentModel, file multipart.File, fileHeader *multipart.FileHeader, db *gorm.DB) error {
	doc.Name = fileHeader.Filename
	db.Create(doc)

	dst, err := os.Create(fmt.Sprintf("docs/%s%s", doc.UUID, filepath.Ext(fileHeader.Filename))) // todo: setup volume on ./docs
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}
	return nil
}
