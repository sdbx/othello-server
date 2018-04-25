package dbs

import (
	"fmt"
	"time"

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
		gorm.Model `json:"-"`
		Winner     string `gorm:"index:idx_battle_winner" json:"winner"`
		Loser      string `gorm:"index:idx_battle_loser" json:"loser"`
		Score      string `json:"score"`
		Moves      string `json:"moves"`
	}

	BattleWithDate struct {
		Battle
		Date time.Time `json:"date"`
	}
)

func GetUserByUserID(UserID string) (User, error) {
	out := User{}
	err := db.Where("user_id = ?", UserID).First(&out).Error

	return out, err
}

func GetUserByName(name string) (User, error) {
	out := User{}
	err := db.Where("name = ?", name).First(&out).Error

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

func (user *User) GetBattles(page int) []BattleWithDate {
	battles := []Battle{}
	db.Where("winner = ? OR loser = ?", user.Name, user.Name).Offset(30 * page).Limit(30).Find(&battles)

	out := []BattleWithDate{}
	for _, battle := range battles {
		out = append(out, BattleWithDate{
			Battle: battle,
			Date:   battle.CreatedAt,
		})
	}

	return out
}

func AddBattle(battle *Battle) error {
	return db.Create(battle).Error
}

func (user *User) GetWinLose() string {
	win := 0
	lose := 0
	db.Model(&Battle{}).Where("winner = ? AND winner <> loser", user.Name).Count(&win)
	db.Model(&Battle{}).Where("loser = ? AND winner <> loser", user.Name).Count(&lose)
	return fmt.Sprintf("%d:%d", win, lose)
}
