package handlers

import (
	"encoding/json"
	"fmt"
	"minpin/url"
	"net/http"
)

type shortenREQ struct {
	URL string `json:"url"`

	Tag string `json:"tag"`
}

type shortenRESP struct {
	Message  string `json:"message"`
	ShortURL string `json:"short_url"`
}

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {

		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req shortenREQ
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	shortURL, err := url.ShortURL(r.Context(), req.Tag, req.URL)
	if err != nil {
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(shortenRESP{
		Message:  "URL shortened successfully",
		ShortURL: fmt.Sprintf("http://localhost:3000/%s", shortURL),
	})
}
