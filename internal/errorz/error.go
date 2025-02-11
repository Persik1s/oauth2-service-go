package errorz

import "errors"

var (
	ErrUserAlreadyExists = errors.New("err user already exists")
	ErrUserNotFound      = errors.New("err user not found")

	ErrPasswordNotValid = errors.New("err password is not valid")

	ErrRoleNotFound = errors.New("err role is not found")

	ErrAuthorization = errors.New("err authorization")
)
