package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// Define a home handler function which writes a byte slice containing
// "Welcome to the Snippetbox Home Page!" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
	 // Check if the current request URL path exactly matches "/". If it doesn't, use
    // the http.NotFound() function to send a 404 response to the client.
    // Importantly, we then return from the handler. If we don't return the handler
    // would keep executing and also write the "Hello from SnippetBox" message.
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

	w.Write([]byte("Welcome to the Snippetbox Home Page!"))
}

 // Add a snippetView handler function.
 func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }

    fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
 }

 // Add a snippetCreate handler function.
 func snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not.
    if r.Method != http.MethodPost {
        w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    w.Write([]byte("Create a new snippet..."))
 }
