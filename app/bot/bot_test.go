package bot

import (
	"context"
	"testing"

	"github.com/jirfag/beepcar-telegram-bot/app/botctx"
	"github.com/jirfag/beepcar-telegram-bot/app/telegram/webhook"
	"github.com/stretchr/testify/assert"
)

func TestProcessWebhook(t *testing.T) {
	err := ProcessWebhook(context.WithValue(context.Background(), botctx.Webhook,
		webhook.RequestPayload{}))
	assert.Nil(t, err)
}
