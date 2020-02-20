package main

import (
	"fmt"
	"log"
	"net/http"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	validationToken := "Validation-Token"
	if len(r.Header.Get(validationToken)) > 0 {
		w.Header().Set(validationToken, r.Header.Get(validationToken))
	}
}

func main() {
	http.HandleFunc("/webhook", webhookHandler)
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
