package user

import "errors"

var (
	ErrorDuplicateUser = errors.New("user already exists")
	ErrorUserPass      = errors.New("user or password not matched")
)
