package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, logErr error) {
	if logErr != nil {
		log.Println(logErr)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// Header setzen – muss vor WriteHeader passieren
	w.Header().Set("Content-Type", "application/json")

	// Payload in JSON serialisieren
	data, err := json.Marshal(payload)
	if err != nil {
		// Serialisierungsfehler → 500 Internal Server Error
		log.Printf("Error marshalling JSON: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Statuscode senden
	w.WriteHeader(code)

	// Daten schreiben und etwaige Schreibfehler protokollieren
	if n, err := w.Write(data); err != nil {
		// n gibt an, wie viele Bytes tatsächlich geschrieben wurden (optional)
		log.Printf("Failed to write JSON response (written %d bytes): %v", n, err)
		// Header sind bereits gesendet – wir können keinen weiteren Statuscode mehr setzen.
	}
}
