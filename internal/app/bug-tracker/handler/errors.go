package handler

import "errors"

var (
	errInvalidJSON               = errors.New("error invalid json")
	errInvalidSignUpData         = errors.New("error invalid sign up data")
	errInvalidSignInData         = errors.New("error invalid sign in data")
	errUserEmailAlreadyExists    = errors.New("error user with such an email already exists")
	errUserUsernameAlreadyExists = errors.New("error user with such a username already exists")
	errInternalServerError       = errors.New("error internal server error")
	errInvalidRefreshToken       = errors.New("error invalid refresh token")
	errTokenDoesNotExist         = errors.New("error token does not exist")
	errUserIsAuthorized          = errors.New("error user is authorized")
	errInvalidAuthHeader         = errors.New("error invalid Authorization header")
	errTokenIsEmpty              = errors.New("error token is empty")
	errUserNotFound              = errors.New("error user is not found")
	errUserDataInvalidType       = errors.New("error user data is of invalid type")
	errEmailIsNotVerified        = errors.New("error email is not verified")
)
