package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jirfag/beepcar-telegram-bot/app/cfg"
)

func CallMethod(method string, payload interface{}) error {
	key := cfg.GetApp().Telegram.APIKey
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", key, method)
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("can't marshal %+v to json: %s", payload, err)
	}

	req, err := http.NewRequest(http.MethodGet, url,
		bytes.NewReader(payloadJSON))
	if err != nil {
		return fmt.Errorf("can't create HTTP request: %s", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("can't execute http request to %q: %s", url, err)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad http status code from %q: %d", url, resp.StatusCode)
	}

	log.Printf("successfully called telegram API method %q with payload %+v",
		method, payload)

	return nil
}
