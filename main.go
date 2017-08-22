package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// TelegramSender ...
type TelegramSender struct {
	ID        int
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// TelegramWebHookPayload is a payload in WebHook message from Telegram API
type TelegramWebHookPayload struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		Date      uint32
		Text      string
		MessageID int `json:"message_id"`
		Chat      struct {
			TelegramSender
			Type string
		}
		From struct {
			TelegramSender
			LanguageCode string `json:"language_code"`
		}
	}
}

func handleUpdateWebHook(r *http.Request) error {
	defer r.Body.Close()
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("can't read request body: %s", err)
	}

	var payload TelegramWebHookPayload
	if err = json.Unmarshal(reqBody, &payload); err != nil {
		return fmt.Errorf("invalid JSON %q in webhook: %s", reqBody, err)
	}

	fmt.Printf("webhook payload is %+v", payload)
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	err := handleUpdateWebHook(r)
	if err != nil {
		log.Printf("request finished with error: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("successfully processed request")
}

func main() {
	var listenAddr string
	flag.StringVar(&listenAddr, "listen-addr", "127.0.0.1:3030", "address to listen")
	flag.Parse()

	log.Printf("listening on %q", listenAddr)

	http.HandleFunc("/btapi/update", handler)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
