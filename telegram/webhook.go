package telegram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type telegramSender struct {
	ID        int
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type telegramWebHookPayload struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		Date      uint32
		Text      string
		MessageID int `json:"message_id"`
		Chat      struct {
			telegramSender
			Type string
		}
		From struct {
			telegramSender
			LanguageCode string `json:"language_code"`
		}
	}
}

func HandleUpdateWebHook(r *http.Request) error {
	defer r.Body.Close()
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("can't read request body: %s", err)
	}

	var payload telegramWebHookPayload
	if err = json.Unmarshal(reqBody, &payload); err != nil {
		return fmt.Errorf("invalid JSON %q in webhook: %s", reqBody, err)
	}

	fmt.Printf("webhook payload is %+v", payload)
	return nil
}
