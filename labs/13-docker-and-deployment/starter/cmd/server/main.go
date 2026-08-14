// Command server is the program you'll containerize in this lab.
// Standard library only — nothing to fetch, nothing to vendor.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

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

// handleOutbound makes a real outbound HTTPS call — used in Exercise 2
// to make the missing-CA-certificates failure observable when the final
// image is built FROM scratch.
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
