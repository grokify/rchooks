package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	validationToken := "Validation-Token"
	if len(r.Header.Get(validationToken)) > 0 {
		w.Header().Set(validationToken, r.Header.Get(validationToken))
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warnf("E_CANNOT_READ_WEBHOOK_BODY [%v]", err.Error())
	} else {
		log.Infof("INCOMING_WEBHOOK_BODY [%v]", string(body))
	}
}

func main() {
	http.HandleFunc("/webhook", webhookHandler)
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
