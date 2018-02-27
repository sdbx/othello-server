package othello

import (
	"math/rand"
	"strconv"
)

type UserStore map[string]string

func (us *UserStore) GenKey(username string) string {
	key := strconv.Itoa(rand.Intn(50000))
	(*us)[key] = username
	return key
}

func (us *UserStore) GetKey(username string) string {
	for key, item := range *us {
		if username == item {
			return key
		}
	}
	return ""
}
