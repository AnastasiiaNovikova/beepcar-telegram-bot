package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jirfag/beepcar-telegram-bot/app/beepcar"
	"github.com/jirfag/beepcar-telegram-bot/app/botctx"
	"github.com/jirfag/beepcar-telegram-bot/app/db"
	"github.com/jirfag/beepcar-telegram-bot/app/models/history"
	"github.com/jirfag/beepcar-telegram-bot/app/models/user"
	"github.com/jirfag/beepcar-telegram-bot/app/telegram/api"
	"github.com/jirfag/beepcar-telegram-bot/app/telegram/webhook"
)

var (
	errInvalidSearchCommand = errors.New("Неверно задана команда, пример: /search Москва Казань")
)

func getUser(ctx context.Context) user.User {
	u := ctx.Value(botctx.User)
	return u.(user.User)
}

func getWebhook(ctx context.Context) webhook.RequestPayload {
	w := ctx.Value(botctx.Webhook)
	return w.(webhook.RequestPayload)
}

func ProcessWebhook(ctx context.Context) error {
	w := getWebhook(ctx)

	user := user.User{
		TelegramID: db.Int64FK(w.Message.From.ID),
	}
	err := user.GetOrCreate()
	if err != nil {
		return fmt.Errorf("can't get/create telegram user: %s", err)
	}

	ctx = context.WithValue(ctx, botctx.User, user)

	payloadJSON, err := json.Marshal(w)
	if err != nil {
		return fmt.Errorf("can't marshal to json webhook: %s", err)
	}

	hw := history.Webhook{
		Payload: string(payloadJSON),
		UserID:  db.Int64FK(int64(user.ID)),
	}
	if err = hw.Save(); err != nil {
		return fmt.Errorf("can't save webhook to history: %s", err)
	}

	if err = processWebhookContent(ctx); err != nil {
		return fmt.Errorf("can't process webhook content: %s", err)
	}

	return nil
}

func processWebhookContent(ctx context.Context) error {
	w := getWebhook(ctx)
	msg := w.Message.Text
	msgFields := strings.Fields(msg)
	if len(msgFields) != 3 || msgFields[0] != "/search" {
		sendToUser(ctx, errInvalidSearchCommand.Error())
		return nil
	}

	fromLocName, toLocName := msgFields[1], msgFields[2]
	fromLocID, err := beepcar.LocationNameToID(ctx, fromLocName)
	if err != nil {
		return fmt.Errorf("can't convert from location name %q to id: %s", fromLocName, err)
	}
	toLocID, err := beepcar.LocationNameToID(ctx, toLocName)
	if err != nil {
		return fmt.Errorf("can't convert to location name %q to id: %s", toLocName, err)
	}

	log.Printf("/search %s %s -> [%d %d]", fromLocName, toLocName,
		fromLocID, toLocID)

	tripIDs, err := beepcar.Search(ctx, fromLocID, toLocID)
	if err != nil {
		return fmt.Errorf("can't search trips from %d to %d: %s",
			fromLocID, toLocID, err)
	}

	tripsLinksMsg := makeTripLinksMsg(tripIDs)
	sendToUser(ctx, tripsLinksMsg)

	return nil
}

func makeTripLinksMsg(tripIDs []int64) string {
	if len(tripIDs) == 0 {
		return "Поездок по маршруту не найдено"
	}

	r := "Поездки:\n"
	for _, tripID := range tripIDs {
		r += fmt.Sprintf("https://beepcar.ru/poezdka/%d\n", tripID)
	}

	return r
}

func sendToUser(ctx context.Context, msg string) {
	w := getWebhook(ctx)
	payload := struct {
		ChatID int64  `json:"chat_id"`
		Text   string `json:"text"`
	}{
		ChatID: w.Message.Chat.ID,
		Text:   msg,
	}

	if err := api.CallMethod("sendMessage", payload); err != nil {
		log.Printf("failed to send message %q to user: %s", msg, err)
	}
}
