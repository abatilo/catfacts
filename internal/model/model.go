package model

import "gorm.io/gorm"

// Target is a single user registration
type Target struct {
	gorm.Model
	PhoneNumber string `gorm:"unique;"`
}
