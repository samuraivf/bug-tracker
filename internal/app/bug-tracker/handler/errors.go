package handler

import "errors"

var (
	errInvalidSignUpData         = errors.New("error invalid sign up data")
	errInvalidSignInData         = errors.New("error invalid sign in data")
	errUserEmailAlreadyExists    = errors.New("error user with such an email already exists")
	errUserUsernameAlreadyExists = errors.New("error user with such a username already exists")
)
