package models

import "gorm.io/gorm"

// User represents the user model
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"not null"`
	Name     string `gorm:"not null"`
}

// UserSession represents the user session model
type UserSession struct {
	gorm.Model
	UserID uint   `gorm:"not null"`
	Token  string `gorm:"not null"`
}

// You can add more fields or customize the model according to your needs
