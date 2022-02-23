package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	validationToken := "Validation-Token"
	if len(r.Header.Get(validationToken)) > 0 {
		w.Header().Set(validationToken, r.Header.Get(validationToken))
		log.Printf("INCOMING_WEBHOOK_VALIDATION_TOKEN [%v]", r.Header.Get(validationToken))
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("E_CANNOT_READ_WEBHOOK_BODY [%v]", err.Error())
	} else {
		log.Printf("INCOMING_WEBHOOK_BODY [%v]", string(body))
	}
}

func main() {
	http.HandleFunc("/webhook", webhookHandler)
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
