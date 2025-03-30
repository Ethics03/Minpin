package handlers

import (
	"context"
	"fmt"
	"minpin/db"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")
	fmt.Println("Requested Short URL:", shortURL)

	longURL, err := db.GetLongURL(context.Background(), shortURL)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, longURL, http.StatusFound)
	fmt.Println("Redirecting to:", longURL)

}
