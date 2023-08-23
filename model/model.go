package model

import (
	"gorm.io/gorm"
)

// User is a struct that represents a user
type User struct {
	ID            uint   `gorm:"primarykey"`
	Name          string `json:"name"`
	Avatar_type   string `json:"avatar_type"`
	Avatar_name   string `json:"avatar_name"`
	Age           int    `json:"age"`
	Year_of_birth int    `json:"year_of_birth"`
	Note          string `json:"note,omitempty"`
	Email         string `json:"email" gorm:"unique"`
	gorm.Model
	// Avatar []byte `json:"-"`
}
