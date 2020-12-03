package main

import (
	"fmt"
	"log"
	"net/http"

	p "github.com/DazWilkin/webmention"
)

const (
	PORT string = "8080"
)

func main() {
	http.HandleFunc("/healthz", p.Healthz)
	http.HandleFunc("/webmention", p.Webmention)
	log.Printf("Listening [:%s]", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil))
}
