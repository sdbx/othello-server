package dbs

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

func init() {
	_init()
}

func _init() {
	tdb, err := gorm.Open("sqlite3", "othello.db")
	if err != nil {
		panic("failed to connect database")
	}
	db = tdb
	migrate()
}

func migrate() {
	db.AutoMigrate(&User{}, &Battle{})
}

func Clear() {
	db.Close()
	err := os.Remove("othello.db")
	if err != nil {
		panic("failed to clear db")
	}
	_init()
}
