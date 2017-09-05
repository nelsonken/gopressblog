package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User table
type User struct {
	gorm.Model
	Name      string `gorm:"not null;type:varchar(20);unique"`
	Nick      string `gorm:"type:varchar(20)"`
	Password  string `gorm:"not null;type:varchar(100)"`
	Avatar    string `gorm:"not null;default empty string"`
	Account   Account
	LastLogin time.Time
	Posts     []Post
	Comments  []Comment
	Message   []Message
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

// Login check in
func (u *User) Login(ORM *gorm.DB, name, password string) error {
	notFound := ORM.First(u, "name = ?", name).RecordNotFound()
	if notFound {
		return LoginError{"用户名不存在", UserNotFound, nil}
	}

	err := u.CheckPassword(password)
	if err != nil {
		return LoginError{"用户名或密码不正确", PasswordError, err}
	}

	if ORM.Model(u).Update("last_login", time.Now()).RowsAffected < 1 {
		return LoginError{"操作失败", DBWriteError, nil}
	}

	return nil
}

// Register register a user
func (u *User) Register(ORM *gorm.DB, name, password string) error {
	notExists := ORM.First(u, "name = ?", name).RecordNotFound()
	if !notExists {
		return RegisterError{"用户名已存在", UserNameExists, nil}
	}

	u.Name = name
	u.LastLogin = time.Now()
	u.GeneratePassword(password)

	if ORM.Save(u).RowsAffected < 1 {
		return RegisterError{"注册失败", DBWriteError, nil}
	}

	return nil
}
