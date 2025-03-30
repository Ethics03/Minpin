package main

import (
	"fmt"
	"log"
	"minpin/db"
	"minpin/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize db: %v\n", err)
	}

	defer db.CloseDB()

	r := chi.NewRouter()

	r.Post("/shorten", handlers.ShortenHandler)
	r.Get("/{shortURL}", handlers.RedirectHandler)

	fmt.Println("Server is running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
