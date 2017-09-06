package models

import (
	"fmt"
)

const (
	// DBWriteError database write failed
	DBWriteError = 0
	// DBReadError database read failed
	DBReadError = -1
	// UserNotFound user name not exists
	UserNotFound = 1
	// PasswordError password error
	PasswordError = 2
	// UserNameExists duplicate username
	UserNameExists = 3
)

// LoginError login err
type LoginError struct {
	Message string
	ErrType int
	Err     error
}

// Error error interface
func (le LoginError) Error() string {
	return fmt.Sprintf("%s", le.Message)
}

// RegisterError reg err
type RegisterError struct {
	Message string
	ErrType int
	Err     error
}

// Error error interface
func (re RegisterError) Error() string {
	return fmt.Sprintf("%s", re.Message)
}

// DBError db error
type DBError struct {
	Message string
	ErrType int
	Err     error
}

// Error error interface
func (de DBError) Error() string {
	return fmt.Sprintf("%s:%v", de.Message, de.Err)
}
