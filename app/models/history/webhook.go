package history

import (
	"database/sql"

	"github.com/jinzhu/gorm"
	"github.com/jirfag/beepcar-telegram-bot/app/cfg"
	"github.com/jirfag/beepcar-telegram-bot/app/db"
	"github.com/jirfag/beepcar-telegram-bot/app/models/user"
)

type Webhook struct {
	gorm.Model

	Payload string
	User    user.User
	UserID  sql.NullInt64 `gorm:"not null;index"`
}

func (w *Webhook) Save() error {
	return db.Get().Create(w).Error
}

func init() {
	if !cfg.IsProduction() {
		db.Get().AutoMigrate(&Webhook{})
	}
}
