package handlers

type URLRequest struct {
	Tag     string `json:"tag"`
	LongURL string `json:"long_url"`
}

type ShortenURL struct {
	ShortURL string `json:"short_url"`
}
