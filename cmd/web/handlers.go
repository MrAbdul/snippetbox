package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox.abdulalsh.com/internal/models"
	"snippetbox.abdulalsh.com/internal/vaildator"
	"strconv"
)

// Define a snippetCreateForm struct to represent the form data and validation
// errors for the form fields. Note that all the struct fields are deliberately
// exported (i.e. start with a capital letter). This is because struct fields
// must be exported in order to be read by the html/template package when
// rendering the template.
type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

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
	data := app.newTemplateData(r)
	// Initialize a new createSnippetForm instance and pass it to the template.
	// Notice how this is also a great opportunity to set any default or
	// 'initial' values for the form --- here we set the initial value for the
	// snippet expiry to 365 days.
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, r, 200, "create.gohtml", data)
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

	// Create an instance of the snippetCreateForm struct containing the values
	// from the form and an empty map for any validation errors.
	form := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}
	// Because the Validator struct is embedded by the snippetCreateForm struct,
	// we can call CheckField() directly on it to execute our validation checks.
	// CheckField() will add the provided key and error message to the
	// FieldErrors map if the check does not evaluate to true. For example, in
	// the first line here we "check that the form.Title field is not blank". In
	// the second, we "check that the form.Title field has a maximum character
	// length of 100" and so on.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	//if there are any errors dump them as plain text and return
	if !form.Valid() {
		//if there are any errors redisply the create.gohtml tmplate passing in the snippetCreateForm instance as dynamic data
		//we use 422 unprocessable entity to indicate that there was a validation error
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.gohtml", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {

		app.serverError(w, r, err)
		return
	}
	//redirect the user to the relevant page fro the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}
