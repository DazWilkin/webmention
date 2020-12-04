package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DazWilkin/webmention"
)

const (
	port string = "8080"
)

func main() {
	http.HandleFunc("/healthz", webmention.Healthz)
	http.HandleFunc("/webmention", webmention.Webmention)
	log.Printf("Listening [:%s]", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
