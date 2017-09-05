package models

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/jinzhu/gorm"
)

// User table
type User struct {
	gorm.Model
	Name     string `gorm:"not null;type:varchar(20);unique"`
	Nick     string `gorm:"type:varchar(20)"`
	Password string `gorm:"not null;type:varchar(100)"`
	Account  Account
	Posts    []Post
	Comments []Comment
	Message  []Message
}

// TableName provide tabel naem to gorm
func (u *User) TableName() string {
	return "users"
}

// GeneratePassword generate password
func (u *User) GeneratePassword(password string) error {
	pb, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(pb)

	return nil
}

// CheckPassword check password
func (u *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	return err
}
