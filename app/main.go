package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jirfag/beepcar-telegram-bot/app/bot"
	"github.com/jirfag/beepcar-telegram-bot/app/telegram/webhook"
)

func main() {
	var listenAddr string
	flag.StringVar(&listenAddr, "listen-addr", ":3030",
		"address to listen")
	flag.Parse()

	log.Printf("listening on %q", listenAddr)

	webhook.SetupWebHookHandlers()
	webhook.SetProcessor(bot.ProcessWebhook)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
