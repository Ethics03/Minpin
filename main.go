package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type shortenREQ struct {
	Url string
	Tag string
}

func main() {
	router := http.NewServeMux()
	//Not using a framework for learning purposes
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	router.HandleFunc("POST /shorten", func(w http.ResponseWriter, r *http.Request) {
		requestData := shortenREQ{}
		err := json.NewDecoder(r.Body).Decode(&requestData) //decodes whatever where is in the request body to json
		if err != nil {
			panic(err)
		}
		w.Write([]byte(requestData.Url))
	})

	server := http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	fmt.Println("Listening on port: 3000")
	server.ListenAndServe()

}
