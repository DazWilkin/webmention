package main

import (
	"fmt"
	"log"
	"net/http"

	p "github.com/DazWilkin/webmention"
)

const (
	port string = "8080"
)

func main() {
	http.HandleFunc("/healthz", p.Healthz)
	http.HandleFunc("/webmention", p.Webmention)
	log.Printf("Listening [:%s]", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
