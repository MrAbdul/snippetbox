package main

import (
	"log"
	"net/http"
)

// this is the home handler which will write a byte slice contiaing the word hello from snippit  box
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox!"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet!"))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet"))
}
func main() {

	//now that we have a handler above (home) we need a router, in go termiology its called servemux
	mux := http.NewServeMux()
	//now that we have our servemux, we can register our handler for the "/" URL pattern
	//Now that the two new routes are up and running, let’s talk a bit of theory.
	//
	//It’s important to know that Go’s servemux has different matching rules depending on whether a route pattern ends
	//with a trailing slash or not.
	//
	//Our two new route patterns — "/snippet/view" and "/snippet/create" — don’t end in a trailing slash. When a pattern
	//doesn’t have a trailing slash, it will only be matched (and the corresponding handler called) when the request URL
	//path exactly matches the pattern in full.
	//
	//When a route pattern ends with a trailing slash — like "/" or "/static/" — it is known as a subtree path pattern.
	//Subtree path patterns are matched (and the corresponding handler called) whenever the start of a request URL path
	//matches the subtree path. If it helps your understanding, you can think of subtree paths as acting a bit like they
	//have a wildcard at the end, like "/**" or "/static/**".
	//
	//This helps explain why the "/" route pattern acts like a catch-all. The pattern essentially means match a single
	//slash, followed by anything (or nothing at all).
	mux.HandleFunc("/{$}", home) // to prevent subtree path patterns from acting like they have a wild card at the end we can append {$} to the end of the pattern so it matches the exact path only. in this case /
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("starting server on :4000")

	// we use the http package to start a new web server, it takes the TCP network address to listen on and the servemux we just created
	err := http.ListenAndServe(":4000", mux)

	//any error returned by the web server is not null and we will log it fatally
	log.Fatal(err)

}
