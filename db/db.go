package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func main() {

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
}

// InsertURL saves a new URL into the database
func InsertURL(tag, longURL, shortURL string) error {
	query := `INSERT INTO urls (tag, long_url, short_url) VALUES ($1, $2, $3) ON CONFLICT (tag) DO NOTHING`
	_, err := Conn.Exec(context.Background(), query, tag, longURL, shortURL)
	return err
}

func GetLongURL(shortURL string) (string, error) {
	var longURL string
	err := Conn.QueryRow(context.Background(), "SELECT long_url FROM urls WHERE short_url=$1", shortURL).Scan(&longURL)
	if err != nil {
		return "", err
	}
	return longURL, nil
}

func CloseDB() {
	Conn.Close(context.Background())
}
