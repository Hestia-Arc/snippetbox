package main

import (
	"fmt"
	"html/template"
	"log"
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

    // Initialize a slice containing the paths to the two files. It's important
    // to note that the file containing our base template must be the *first*
    // file in the slice.
    files := []string{
        "./ui/html/base.tmpl",
        "./ui/html/partials/nav.tmpl",
        "./ui/html/pages/home.tmpl",
    }

    // Use the template.ParseFiles() function to read the template file into a
    // template set. If there's an error, we log the detailed error message and use
    // the http.Error() function to send a generic 500 Internal Server Error
    // response to the user.
    // ts, err := template.ParseFiles("./ui/html/pages/home.tmpl")
    ts, err := template.ParseFiles(files...)
    if err != nil {
        log.Println(err.Error())
        http.Error(w, "Internal Server Error", 500)
        return
    }

    // We then use the Execute() method on the template set to write the
    // template content as the response body. The last parameter to Execute()
    // represents any dynamic data that we want to pass in, which for now we'll
    // leave as nil.
    // err = ts.Execute(w, nil)
    err = ts.ExecuteTemplate(w, "base", nil)
    if err != nil {
        log.Println(err.Error())
        http.Error(w, "Internal Server Error", 500)
    }

	// w.Write([]byte("Welcome to the Snippetbox Home Page!"))
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
