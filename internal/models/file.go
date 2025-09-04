package models

import (
	"mime/multipart"

	"gorm.io/gorm"
)

// File Database Schema for File
type File struct {
	gorm.Model
	UserID   int64  `gorm:"column:user_id"`
	FileName string `gorm:"column:file_name"`
	FileURL  string `gorm:"column:file_url"`
}

// Schema for Upload File
type UploadFile struct {
	FileName    string
	FileData    multipart.File
	Size        int64
	ContentType string
	FilePath    string
}
