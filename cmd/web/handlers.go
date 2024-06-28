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
	data := TemplateData{
		Snippets: s,
	}

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

	data := TemplateData{Snippet: s}
	app.render(w, r, 200, "view.gohtml", data)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet"))
}
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	//dummy data
	title := "0snail"
	content := "0 snail \n climb mount fuji, \n But slowly,slowly! \n\n-kobayashi issa"
	expires := 7
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	//redirect the user to the relevant page fro the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}
