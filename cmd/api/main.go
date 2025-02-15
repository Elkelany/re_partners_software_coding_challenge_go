package main

import (
	"log"

	"re_partners_software_coding_challenge_go/cmd/api/internal/http"
)

func main() {
	// Create a new instance of http.Server
	srv := http.NewServer()

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
