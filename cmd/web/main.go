package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"snippetbox.abdulalsh.com/internal/models"

	_ "github.com/go-sql-driver/mysql" //imported for affect
)

// we define an application struct to hole the application wide dependencies for the web application
// for now we will only include the structred logger but we will be adding more later on
type application struct {
	logger *slog.Logger
	//we add a snippets field to the application struct to make the Snippet model available to our handlers
	snippets *models.SnippetModel
	//the template cache
	templateCache map[string]*template.Template
}

func main() {

	//we define a new command line flag with the name addr with a default value of :4000
	//and some short helo text explaining what the flag controls
	addr := flag.String("addr", ":4000", "http service address")
	//we defined a new command-line flag for mysql dsn string
	dsn := flag.String("dsn", "web2:TOOR@/snippetbox?parseTime=true", "MySQL datasource name")
	//we have to call flag.Parse() before any variable of the command line is used, if any errors are encontred during the parsing
	//the application will be terminated
	flag.Parse()

	//lets add a structred logger to our applicattion
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	//we also defer a call to db.Close() so that the connection pool is closed before the main function exists
	defer db.Close()
	// Initialize a new instance of our application struct, containing the
	// dependencies (for now, just the structured logger).

	cache, err := newTemplateCache()
	if err != nil {
		logger.Error("problem initializing template cache", err)
		os.Exit(1)
	}

	app := &application{
		logger: logger,
		//we init a models.snippetmodel instance with the connection pool and add it to the application depencies
		snippets:      &models.SnippetModel{DB: db},
		templateCache: cache,
	}
	//now that we have a handler above (home) we need a router, in go termiology its called servemux

	//the value returned from the flag.String is a pointer to the flag value, not the value itself.
	//so we need to defrence it with *
	logger.Info("starting server", "addr", *addr)

	// we use the http package to start a new web server, it takes the TCP network address to listen on and the servemux we just created
	// and we defrence it here as well
	// Call the new app.routes() method to get the servemux containing our routes,
	// and pass that to http.ListenAndServe().
	err = http.ListenAndServe(*addr, app.routes())

	//any error returned by the web server is not null and we will log it fatally
	logger.Error(err.Error())
	os.Exit(1)

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	//the sql.open doesnt actually create anny connections, it just initilize the pool for future use
	//the connections are established lazily, so to verify that everything is setup correctly we use the db.ping
	//this will create a connection and check for errors, if we get an error we call db.close to close the pool and return the error
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
