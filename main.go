package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func handleUpdateWebHook(r *http.Request) error {
	defer r.Body.Close()
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("can't read request body: %s", err)
	}

	fmt.Printf("request body is %q", reqBody)
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	err := handleUpdateWebHook(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/btapi/update", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
