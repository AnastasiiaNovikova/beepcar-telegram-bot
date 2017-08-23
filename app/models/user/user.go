package user

import (
	"database/sql"

	"github.com/jinzhu/gorm"
	"github.com/jirfag/beepcar-telegram-bot/app/db"
)

type User struct {
	gorm.Model

	TelegramID sql.NullInt64 `gorm:"not null;unique_index"`
	FirstName  string
	SecondName string
}

func (u *User) GetOrCreate() error {
	// XXX: RC here
	return db.Get().FirstOrCreate(u, User{TelegramID: u.TelegramID}).Error
}

func init() {
	//	if !cfg.IsProduction() {
	db.Get().AutoMigrate(&User{})
	//	}
}
