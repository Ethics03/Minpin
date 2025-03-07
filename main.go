package main

import (
	"encoding/json"
	"fmt"
	"log"
	"minpin/db"
	"minpin/url"
	"net/http"
)

// Request structure for shortening URLs
type shortenREQ struct {
	URL string `json:"url"`
	Tag string `json:"tag"`
}

// Response structure
type shortenRESP struct {
	Message  string `json:"message"`
	ShortURL string `json:"short_url"`
}

func main() {
	// Initialize the database
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v\n", err)
	}
	defer db.CloseDB()

	router := http.NewServeMux()

	// Root Route
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Write([]byte("URL Shortener API is running"))
	})

	// Shorten URL Route
	router.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var req shortenREQ
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Generate short URL and store it in DB
		shortURL, err := url.ShortURL(r.Context(), req.Tag, req.URL)
		if err != nil {
			http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
			return
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(shortenRESP{
			Message:  "URL shortened successfully",
			ShortURL: shortURL,
		})
	})

	// Resolve URL Route
	router.HandleFunc("/resolve", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract the short URL from query params
		shortURL := r.URL.Query().Get("short")
		if shortURL == "" {
			http.Error(w, "Short URL parameter missing", http.StatusBadRequest)
			return
		}

		// Retrieve original URL
		longURL, err := url.ResolveURL(r.Context(), shortURL)
		if err != nil {
			http.Error(w, "Short URL not found", http.StatusNotFound)
			return
		}

		// Redirect to original URL
		http.Redirect(w, r, longURL, http.StatusFound)
	})

	// Start the server
	server := http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	fmt.Println("🚀 Server is running on port 3000")
	log.Fatal(server.ListenAndServe())
}
