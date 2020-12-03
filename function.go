package p

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"cloud.google.com/go/firestore"
)

// Obtain ProjectID from the environment (provided by Cloud Functions)
var projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

var client *firestore.Client

// Mention is a type that includes metadata associated with a mention
type Mention struct {
	Datetime time.Time
}

// Firestore client initialized here; once per Instance (not per Function)
func init() {
	var err error

	ctx := context.Background()
	client, err = firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal("firestore.Client: %v", err)
	}
}

// Webmention is a handler that receives `POST`ed webmentions
// Expects headers for Source(Mentioner)|Target(Mentioned)
func Webmention(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	log.Printf("[Webmention] Content-Type: %s", contentType)

	host := r.Header.Get("Host")
	log.Printf("[Webmention] Host: %s", host)

	// Where is the mention?
	sourceURL := r.Header.Get("Source")
	log.Printf("[Webmention] source: %s", sourceURL)

	if sourceURL == "" {
		log.Print("[Webmention] Required header (`Source`) omitted")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request: Source header required"))
		return
	}

	// Parse source URL
	source, err := url.Parse(sourceURL)
	if err != nil {
		log.Print("[Webmention] Unable to parse header (`Source`)")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request: Unable to parse Source header"))
		return
	}

	// What was mentioned?
	targetURL := r.Header.Get("Target")
	log.Printf("[Webmention] target: %s", targetURL)

	if targetURL == "" {
		log.Print("[Webmention] Required header (`Target`) omitted")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request: Target header required"))
		return
	}

	// Parse target URL
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Print("[Webmention] Unable to parse header (`Target`)")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request: Unable to parse Target header"))
		return
	}

	// Order by:
	// + what was mentioned (i.e. target)
	// + what is the mention (i.e. source)
	// `.Path` values retain `/` prefix so it need not be duplicated
	path := fmt.Sprintf("Mentions/%s%s/%s%s", target.Host, target.Path, source.Host, source.Path)
	log.Printf("[Webmention] Document: %s", path)

	mention := client.Doc(path)

	// Uses the handler's context for Firestore
	_, err = mention.Create(r.Context(), Mention{
		Datetime: time.Now(),
	})
	if err != nil {
		log.Printf("[Webmention] Unable to write document: %s\n%v", path, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error: Unable to record"))
		return
	}

	fmt.Fprintln(w, "ok")
}

// Healthz is a handler that returns the function's health (oK)
func Healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}
