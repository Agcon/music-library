package models

import "gorm.io/gorm"

type Song struct {
	ID       uint32 `gorm:"primaryKey;autoIncrement" json:"id"`
	Title    string `gorm:"not null" json:"title" binding:"required"`
	Group    string `gorm:"not null" json:"group" binding:"required"`
	Text     string `gorm:"not null" json:"text,omitempty"`
	FilePath string `gorm:"not null" json:"filePath,omitempty"`
}

func MigrateSong(db *gorm.DB) error {
	return db.AutoMigrate(&Song{})
}
