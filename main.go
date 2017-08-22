package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jirfag/beepcar-telegram-bot/telegram"
)

func handler(w http.ResponseWriter, r *http.Request) {
	err := telegram.HandleUpdateWebHook(r)
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
