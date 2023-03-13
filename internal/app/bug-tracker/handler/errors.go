package handler

import "errors"

var (
	errInvalidSignUpData = errors.New("error invalid sign up data")
	errInvalidSignInData = errors.New("error invalid sign in data")
)