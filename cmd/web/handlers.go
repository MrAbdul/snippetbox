package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox.abdulalsh.com/internal/models"
	"strconv"
)

// this is the home handler which will write a byte slice contiaing the word hello from snippit  box
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//	panic("oops! something went wrong") // Deliberate panic
	//Unfortunately, all we get is an empty response due to Go closing the underlying HTTP connection following the panic.
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	//for _, snippet := range snippets {
	//	fmt.Fprintf(w, "%+v\n", snippet)
	//}

	//2.6 you must ensure that header map contains all the headers you want before calling w.writeheader or w.write
	//2.8 we can use the template here
	//2.8 part 2 we can parse multiple template filesm
	//		note that the file containg the base template must be first in the slice
	//5.3 now we see tha payoff of our work
	//5.5 we will call newTemplateData helper to get the templatedata struct containg the default data, which for now is just the current year
	data := app.newTemplateData(r)
	data.Snippets = s

	app.render(w, r, 200, "home.gohtml", data)

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	//we can retrive a wild card id like so by refereing to its wildcard slug name
	id := r.PathValue("id")
	//since this will be an untrusted user input we should validate it to make sure its sensible before we use it
	//for this case, we need to make sure it is a positive integer
	idint, err := strconv.Atoi(id)
	if err != nil || idint < 1 {
		http.NotFound(w, r)
		return
	}

	s, err := app.snippets.Get(idint)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	//5.5 we will call newTemplateData helper to get the templatedata struct containg the default data, which for now is just the current year
	data := app.newTemplateData(r)
	data.Snippet = s
	app.render(w, r, 200, "view.gohtml", data)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, 200, "create.gohtml", app.newTemplateData(r))
}
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// First we call r.ParseForm() which adds any data in POST request bodies
	// to the r.PostForm map. This also works in the same way for PUT and PATCH
	// requests. If there are any errors, we use our app.ClientError() helper to
	// send a 400 Bad Request response to the user.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use the r.PostForm.Get() method to retrieve the title and content
	// from the r.PostForm map.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	// The r.PostForm.Get() method always returns the form data as a *string*.
	// However, we're expecting our expires value to be a number, and want to
	// represent it in our Go code as an integer. So we need to manually covert
	// the form data to an integer using strconv.Atoi(), and we send a 400 Bad
	// Request response if the conversion fails.
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	//redirect the user to the relevant page fro the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}
