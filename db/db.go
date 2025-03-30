package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var pool *pgxpool.Pool

func InitDB() error {

	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env not set")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DB URL not set")
	}

	fmt.Println("Connecting to database:", dbURL)

	pool, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("database is not reachable: %w", err)
	}

	fmt.Println("Database connected successfully")
	return nil
}

func InsertURL(ctx context.Context, tag, longURL, shortURL string) error {
	query := `INSERT INTO urls (tag, long_url, short_url) VALUES ($1, $2, $3) ON CONFLICT (tag) DO NOTHING`
	_, err := pool.Exec(ctx, query, tag, longURL, shortURL)
	return err
}

func GetLongURL(ctx context.Context, shortURL string) (string, error) {
	var longURL string
	err := pool.QueryRow(ctx, "SELECT long_url FROM urls WHERE short_url=$1", shortURL).Scan(&longURL)
	if err != nil {
		return "", err
	}
	return longURL, nil
}

func CloseDB() {
	if pool != nil {
		pool.Close()
		fmt.Println("DB CLOSED")
	}

}

func ShortExists(ctx context.Context, shortURL string) (bool, error) {
	var exists bool
	err := pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM urls WHERE short_url=$1)", shortURL).Scan(&exists)
	return exists, err
}
