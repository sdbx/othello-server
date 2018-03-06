package db

import (
	"math/rand"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sdbx/othello-server/othello/models"
)

type (
	DBUserStore map[string]*models.User
)

var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func genKey(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (us *DBUserStore) Register(username string) string {
	if len(username) == 0 {
		return ""
	}

	for _, item := range *us {
		if item.Name == username {
			return ""
		}
	}

	key := genKey(10)
	(*us)[key] = &models.User{
		Name:   username,
		Secret: key,
	}
	return key
}

func (us *DBUserStore) GetUserByID(username string) *models.User {
	for _, item := range *us {
		if username == item.Name {
			return item
		}
	}
	return nil
}

func (us *DBUserStore) GetUserBySecret(secret string) *models.User {
	if user, ok := (*us)[secret]; !ok {
		return nil
	} else {
		return user
	}
}
