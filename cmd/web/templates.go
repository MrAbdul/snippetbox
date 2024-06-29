package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"snippetbox.abdulalsh.com/internal/models"
	"time"
)

type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
	Form        any
	Flash       string
}

// we create a newTemplateData() helper which will return a pointer to a templatedata struct initilized with the current year,
// note that we are not using the http.request param, but we will do later
func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
		// Use the PopString() method to retrieve the value for the "flash" key.
		// PopString() also deletes the key and value from the session data, so it
		// acts like a one-time fetch. If there is no matching key in the session
		// data this will return the empty string.
		// Add the flash message to the template data, if one exists.
		Flash: app.sessionManager.PopString(r.Context(), "flash"),
	}
}
func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get a slice of all filepaths that
	// match the pattern "./ui/html/pages/*.tmpl". This will essentially gives
	// us a slice of all the filepaths for our application 'page' templates
	// like: [ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
	pages, err := filepath.Glob("./ui/html/pages/*.gohtml")
	if err != nil {
		return nil, err
	}
	// Loop through the page filepaths one-by-one.
	for _, page := range pages {
		// Extract the file name (like 'home.gohtml') from the full filepath
		// and assign it to the name variable.
		name := filepath.Base(page)
		// Parse the base template file into a template set.
		//the func map must be registered with the template set before we call the parse files
		//this means we have to use templates.New to create an empty template set, use the funcs() method to register the template func map then parse the files as normal

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.gohtml")
		if err != nil {
			return nil, err
		}
		// Call ParseGlob() *on this template set* to add any partials.
		ts, err = ts.ParseGlob("./ui/html/partials/*.gohtml")
		if err != nil {
			return nil, err
		}
		// Create a slice containing the filepaths for our base template, any
		// partials and the page.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// Add the template set to the map, using the name of the page
		// (like 'home.gohtml') as the key.
		// Add the template set to the map as normal...

		cache[name] = ts
	}
	// Return the map.
	return cache, nil

}

// creating custome functions for templates
func humanDate(t time.Time) string {
	return t.Format("02 01 2006 at 15:04")
}

// we will init a template.FuncMap object and store it in a global variable. this is essentially a string-keyed map which acts as a lookup
// between the names of our custom template functions and the function itself
var functions = template.FuncMap{
	"humanDate": humanDate,
}
