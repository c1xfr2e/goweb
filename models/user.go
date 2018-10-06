package models

import (
	"github.com/jinzhu/gorm"
)

// User for authorized user entity
type User struct {
	ID       uint   `gorm:"primary_key;auto_increment"`
	Email    string `gorm:"column:email;not null"`
	Password string `gorm:"column:password;not null"`
}

// TableName defines table name
func (User) TableName() string {
	return "users"
}

// GetUserByEmail gets user from DB by email
func GetUserByEmail(db *gorm.DB, email string) *User {
	ret := new(User)
	err := db.Where("email = ?", email).First(ret).Error
	if err != nil {
		return nil
	}
	return ret
}

// GetUserByID get user from DB by ID
func GetUserByID(db *gorm.DB, id uint) *User {
	ret := new(User)
	err := db.Where("id = ?", id).First(ret).Error
	if err != nil {
		return nil
	}
	return ret
}
