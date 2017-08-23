package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jirfag/beepcar-telegram-bot/app/botctx"
)

const apiPath = "/btapi/update"

type webhookProcessor func(ctx context.Context) error

var processor webhookProcessor

func SetProcessor(p webhookProcessor) {
	processor = p
}

type sender struct {
	ID        int64
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

type RequestPayload struct {
	UpdateID int64 `json:"update_id"`
	Message  struct {
		Date      uint32
		Text      string
		MessageID int `json:"message_id"`
		Chat      struct {
			sender
			Type string
		}
		From struct {
			sender
			LanguageCode string `json:"language_code"`
		}
	}
}

func handleUpdateRequest(r *http.Request) error {
	defer r.Body.Close()
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("can't read request body: %s", err)
	}

	var payload RequestPayload
	if err = json.Unmarshal(reqBody, &payload); err != nil {
		return fmt.Errorf("invalid JSON %q in webhook: %s", reqBody, err)
	}

	ctx := context.WithValue(context.Background(), botctx.Webhook, payload)

	if processor != nil {
		if err = processor(ctx); err != nil {
			return fmt.Errorf("can't process webhook payload: %s", err)
		}
	}

	fmt.Printf("webhook payload is %+v", payload)
	return nil
}

func handleHTTPRequest(w http.ResponseWriter, r *http.Request) {
	err := handleUpdateRequest(r)
	if err != nil {
		log.Printf("request finished with error: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("successfully processed request")
}

func SetupWebHookHandlers() {
	http.HandleFunc(apiPath, handleHTTPRequest)
}
