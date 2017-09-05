package models

import (
	"fmt"
)

const (
	DBWriteError   = 0
	DBReadError    = -1
	UserNotFound   = 1
	PasswordError  = 2
	UserNameExists = 3
)

// LoginError
type LoginError struct {
	Message string
	ErrType int
	Err     error
}

// Error error interface
func (le LoginError) Error() string {
	return fmt.Sprintf("%s:%v", le.Message, le.Err)
}

// LoginError
type RegisterError struct {
	Message string
	ErrType int
	Err     error
}

// Error error interface
func (re RegisterError) Error() string {
	return fmt.Sprintf("%s:%v", re.Message, re.Err)
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
