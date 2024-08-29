package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
	Blogs    []Blog `gorm:"foreignKey:UserID" json:"blogs"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
