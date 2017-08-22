package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
		MessageID string `json:"message_id"`
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
	var payload TelegramWebHookPayload
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		return fmt.Errorf("invalid JSON in webhook: %s", err)
	}

	fmt.Printf("webhook payload is %v", payload)
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
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

	http.HandleFunc("/btapi/update", handler)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
