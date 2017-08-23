package bot

import (
	"encoding/json"
	"fmt"

	"github.com/jirfag/beepcar-telegram-bot/app/db"
	"github.com/jirfag/beepcar-telegram-bot/app/models/history"
	"github.com/jirfag/beepcar-telegram-bot/app/models/user"
	"github.com/jirfag/beepcar-telegram-bot/app/telegram/webhook"
)

func ProcessWebhook(r webhook.RequestPayload) error {
	user := user.User{
		TelegramID: db.Int64FK(r.Message.From.ID),
	}
	err := user.GetOrCreate()
	if err != nil {
		return fmt.Errorf("can't get/create telegram user: %s", err)
	}

	payloadJSON, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("can't marshal to json webhook: %s", err)
	}

	w := history.Webhook{
		Payload: string(payloadJSON),
		UserID:  db.Int64FK(int64(user.ID)),
	}
	if err = w.Save(); err != nil {
		return fmt.Errorf("can't save webhook to history: %s", err)
	}

	if err = processWebhookContent(r); err != nil {
		return fmt.Errorf("can't process webhook content: %s", err)
	}

	return nil
}

func processWebhookContent(r webhook.RequestPayload) error {
	return nil
}
