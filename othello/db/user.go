package db

import (
	"math/rand"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sdbx/othello-server/othello/models"
)

type (
	DBUserStore struct {
	}
)

var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()_")

func genKey(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (us *DBUserStore) Register(username string) string {
	key := genKey(10)
	return key
}

func (us *DBUserStore) GetUserByName(username string) *models.User {

	return nil
}

func (us *DBUserStore) GetUserBySecret(secret string) *models.User {
	return nil
}
