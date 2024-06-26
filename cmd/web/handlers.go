package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// this is the home handler which will write a byte slice contiaing the word hello from snippit  box
func home(w http.ResponseWriter, r *http.Request) {
	//2.6 you must ensure that header map contains all the headers you want before calling w.writeheader or w.write
	//2.8 we can use the template here
	files, err := template.ParseFiles("./ui/html/pages/home.gohtml")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "internal server Error", http.StatusInternalServerError)
		return
	}
	//now that we have the file opened we can execute it
	err = files.Execute(w, nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "internal server Error", http.StatusInternalServerError)
	}

}

func snippetView(w http.ResponseWriter, r *http.Request) {
	//we can retrive a wild card id like so by refereing to its wildcard slug name
	id := r.PathValue("id")
	//since this will be an untrusted user input we should validate it to make sure its sensible before we use it
	//for this case, we need to make sure it is a positive integer
	idint, err := strconv.Atoi(id)
	if err != nil || idint < 1 {
		http.NotFound(w, r)
		return
	}
	//we use fmt.Sprintf() to interpolate the id value with a message, then write it as http response
	//since fmt.Fprintf takes a io.writer, we can shorten the following
	//msg := fmt.Sprintf("Display a specifc snippit with ID %d", idint)
	//w.Write([]byte(msg))
	//	to
	fmt.Fprintf(w, "Display a specifc snippit with ID %d", idint)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet"))
}
func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("save a new snippet"))

}