package othello

import (
	"math/rand"
	"strconv"
)

type UserStatus uint

const (
	None UserStatus = iota + 1
	InRoom
)

type User struct {
	Name   string
	Secret string
	Status UserStatus
}

type UserStore map[string]*User

func (us *UserStore) Register(username string) string {
	if len(username) == 0 {
		return ""
	}

	for _, item := range *us {
		if item.Name == username {
			return ""
		}
	}

	key := strconv.Itoa(rand.Intn(50000))
	(*us)[key] = &User{
		Name:   username,
		Secret: key,
	}
	return key
}

func (us *UserStore) GetUserByName(username string) *User {
	for _, item := range *us {
		if username == item.Name {
			return item
		}
	}
	return nil
}

func (us *UserStore) GetUserBySecret(secret string) *User {
	if user, ok := (*us)[secret]; !ok {
		return nil
	} else {
		return user
	}
}
