// Command server is a small HTTP API used as the running example for
// Topic 13 (Docker & Deployment). It's a simplified stand-in for the
// REST service built in Topic 10 — enough shape to be worth
// containerizing, deliberately kept to the standard library only, so
// the image it produces has no third-party dependencies to reason
// about.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// item mirrors the shape of Topic 10's domain.Item, trimmed down.
type item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var items = []item{
	{ID: "1", Name: "widget"},
	{ID: "2", Name: "gadget"},
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// DATABASE_URL isn't wired to a real driver here — this stand-in
	// doesn't need one to demonstrate containerizing a Go service. It's
	// read and logged so docker-compose's environment injection (see
	// docker-compose.yml) has something concrete to point at. A real
	// service would pass this to sql.Open with a Postgres driver.
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		log.Printf("configured with DATABASE_URL=%s", dsn)
	} else {
		log.Println("no DATABASE_URL set — running without a database")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handleHealthz)
	mux.HandleFunc("/items", handleItems)
	mux.HandleFunc("/outbound", handleOutbound)

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func handleItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// handleOutbound makes a real outbound HTTPS call. Its only job is to
// make the "FROM scratch has no CA certificates" gotcha observable
// rather than theoretical: hit this endpoint from a container built on
// scratch without certs copied in, and the request fails with a TLS
// certificate-verification error. Copy /etc/ssl/certs/ca-certificates.crt
// from the builder stage and it starts working.
func handleOutbound(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://api.github.com")
	if err != nil {
		http.Error(w, "outbound call failed: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	w.Write([]byte("outbound call succeeded, status: " + resp.Status))
}
