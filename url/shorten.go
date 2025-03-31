package url

// for shortening the url
import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"minpin/db"
	"strings"
)

func ShortenURL(ctx context.Context) (string, error) {
	const shortl = 6

	for {
		bytes := make([]byte, shortl)
		if _, err := rand.Read(bytes); err != nil {
			return "", err
		}

		shortURL := base64.URLEncoding.EncodeToString(bytes)

		shortURL = strings.TrimRight(shortURL, "=")[:shortl]

		exists, err := db.ShortExists(ctx, shortURL)
		if err != nil {
			return "", err
		}
		if !exists {
			return shortURL, nil
		}
	}
}

func ShortURL(ctx context.Context, tag, longURL string) (string, error) {
	shorturl, err := ShortenURL(ctx)
	if err != nil {
		return "", err
	}

	err = db.InsertURL(ctx, tag, longURL, shorturl)

	if err != nil {
		return "", err
	}

	return shorturl, nil
}

func ResolveURL(ctx context.Context, shortURL string) (string, error) {
	longURL, err := db.GetLongURL(ctx, shortURL)
	if err != nil {
		return "", errors.New("URL not found")
	}
	return longURL, nil
}
