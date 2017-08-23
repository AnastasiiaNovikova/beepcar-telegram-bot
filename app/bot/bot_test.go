package bot

import (
	"testing"

	"github.com/jirfag/beepcar-telegram-bot/app/telegram/webhook"
	"github.com/stretchr/testify/assert"
)

func TestProcessWebhook(t *testing.T) {
	err := ProcessWebhook(webhook.RequestPayload{})
	assert.Nil(t, err)
}
