package dbs

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/sdbx/othello-server/othello/utils"
)

type (
	User struct {
		gorm.Model
		Name    string `grom:"unique"`
		Secret  string `gorm:"unique"`
		Profile string
		UserID  string `gorm:"unique"`
	}

	Battle struct {
		gorm.Model
		Winner string `gorm:"index:idx_battle_winner"`
		Loser  string `gorm:"index:idx_battle_loser"`
		Score  string
		Moves  string
	}
)

func GetUserByUserID(UserID string) (User, error) {
	out := User{}
	err := db.Where("user_id = ?", UserID).First(&out).Error

	return out, err
}

func GetUserBySecret(secret string) (User, error) {
	out := User{}
	err := db.Where("secret = ?", secret).First(&out).Error

	return out, err
}

func secretTest(secret string) bool {
	count := 0
	users := []User{}
	err := db.Where("secret = ?", secret).Find(&users).Count(&count).Error
	if err != nil {
		return false
	}

	return count == 0
}

func nameTest(name string) bool {
	count := 0
	users := []User{}
	err := db.Where("name = ?", name).Find(&users).Count(&count).Error
	if err != nil {
		return false
	}

	return count == 0
}

func AddUser(user *User) error {
	d := 0
	for !nameTest(user.Name) {
		user.Name = fmt.Sprintf("%s%d", user.Name, d)
		d++
	}

	user.Secret = utils.GenKey()
	for !secretTest(user.Secret) {
		user.Secret = utils.GenKey()
	}

	return db.Create(user).Error
}

func (user *User) GetBattles() []Battle {
	battles := []Battle{}
	db.Where("winner = ? OR loser = ?", user.Name, user.Name).Find(&battles)
	return battles
}

func AddBattle(battle *Battle) error {
	return db.Create(battle).Error
}
