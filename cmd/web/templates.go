package main

import (
	"html/template"
	"path/filepath"
	"snippetbox.abdulalsh.com/internal/models"
)

type TemplateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
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
		ts, err := template.ParseFiles("./ui/html/base.gohtml")
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
