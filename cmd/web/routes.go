package main

import (
	"github.com/justinas/alice"
	"net/http"
	"snippetbox.abdulalsh.com/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	//
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
	// Create a new middleware chain containing the middleware specific to our
	// dynamic application routes. For now, this chain will only contain the
	// LoadAndSave session middleware but we'll add more to it later.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	//chapter2.5 method based routing, we can restrict a path to a specific method by prefix the route pattern with the necessery http method
	// Update these routes to use the new dynamic middleware chain followed by
	// the appropriate handler function. Note that because the alice ThenFunc()
	// method returns a http.Handler (rather than a http.HandlerFunc) we also
	// need to switch to registering the route using the mux.Handle() method.
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))                      // to prevent subtree path patterns from acting like they have a wild card at the end we can append {$} to the end of the pattern so it matches the exact path only. in this case /
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView)) //lets include a wildcard segment to select a specific id
	//Notes on wildcard precedence and conflict:
	//if an overlap occurs for example "/post/edit" and "/post/{id}" the first one is a valid match for both patterns
	//the rule for this is succinct: the most specific route pattern wins:
	//go defines a pattern as more specific than another if it matches only a subset of requests that the other pattern matches
	//the /post/edit only matches requests with the exact path /post/edit, whereas the pattern /post/{id} matches requests with
	// /post/edit, /post/123, /post/abc and many more, therefore /post/edit is the more specifc route pattern and will take precedent.
	//1. a side effect of this is that you can register patterns any any order and it wont change how the servmux behaves

	//2.if an edge case occures where two overlapping route patterns arent obvously more specific than the other,
	//	for example "/post/new/{id} and /post/{author}/latest overlap
	//	because they both match the /post/new/latest but its not clear which should take precedence
	// 	go's servermuux considers this as pattern conflict and will panic at runtime when initializing the orutes

	// Swap the route declarations to use the application struct's methods as the
	// handler functions.

	// ch2.5 note that we can create routes that have the same pattern but diffrent HTTP methods

	//2.9 we create a file server to serve files out of the "./ui/static"directory
	fileServer := http.FileServerFS(ui.Files)
	// now we use mux.handle function to register the file server as handler for all url paths that start with /static/
	// for matching paths, we strip the "/static" prefix before the request reaches the file server
	//we don't need the session handling middleware for the static files
	mux.Handle("GET /static/", fileServer)

	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignUp))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignUpPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requiresAuthentication)
	//protected routes
	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	// ch2.5 lets add a post only route and handler
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	//we can use justinas alice to compose the  handlers
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
