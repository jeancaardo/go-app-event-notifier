package users

import "errors"

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrUserEmailAlreadyExists = errors.New("user email already exists")
	ErrOnStoreUser            = errors.New("error on store user")
	ErrOnUpdateUser           = errors.New("error on update user")
	ErrOnDeleteUser           = errors.New("error on delete user")
)
