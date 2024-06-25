package main

import (
	"log"
	"net/http"
)

// this is the home handler which will write a byte slice contiaing the word hello from snippit  box
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox!"))
}

func main() {

	//now that we have a handler above (home) we need a router, in go termiology its called servemux
	mux := http.NewServeMux()
	//now that we have our servemux, we can register our handler for the "/" URL pattern
	mux.HandleFunc("/", home)

	log.Print("starting server on :4000")

	// we use the http package to start a new web server, it takes the TCP network address to listen on and the servemux we just created
	err := http.ListenAndServe(":4000", mux)

	//any error returned by the web server is not null and we will log it fatally
	log.Fatal(err)

}
