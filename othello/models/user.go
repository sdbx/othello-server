package models

type User struct {
	Name   string
	Secret string
}

type UserStore interface {
	Register(userID string) string
	GetUserByID(uesrID string) *User
	GetUserBySecret(secret string) *User
}
