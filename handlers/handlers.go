package handlers

import (
	"context"
	"encoding/json"
	"minpin/db"
	"net/http"
)

type URLRequest struct {
	Tag     string `json:"tag"`
	LongURL string `json:"long_url"`
}

type ShortURL struct {
	ShortURL string `json:"short_url"`
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req URLRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.LongURL == "" {
		http.Error(w, "long_url is required", http.StatusBadRequest)
		return
	}

	shortUrl := req.Tag

	if shortUrl == "" {
		shortUrl = generateShortURL()
	}

	query := `INSERT INTO urls (tag, long_url, short_url) VALUES ($1, $2, $3) ON CONFLICT (tag) DO NOTHING RETURNING short_url`
	var storeShortURL string

	err = db.Conn.QueryRow(context.Background(), query, shortURL, req.LongURL, shortUrl).Scan(&storeShortURL)

}
