package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/grokify/rchooks"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if token := strings.TrimSpace(r.Header.Get(rchooks.HeaderValidationToken)); len(token) > 0 {
		w.Header().Set(rchooks.HeaderValidationToken, token)
		log.Printf("INCOMING_WEBHOOK_VALIDATION_TOKEN [%v]", token)
		return
	}

	body, err := io.ReadAll(r.Body)
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
