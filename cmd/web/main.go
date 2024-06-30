package main

import (
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"snippetbox.abdulalsh.com/internal/models"
	"time"

	_ "github.com/go-sql-driver/mysql" //imported for affect
)

// we define an application struct to hole the application wide dependencies for the web application
// for now we will only include the structred logger but we will be adding more later on
type application struct {
	logger *slog.Logger
	//we add a snippets field to the application struct to make the Snippet model available to our handlers
	snippets *models.SnippetModel
	//the template cache
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
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
	// Use the scs.New() function to initialize a new session manager. Then we
	// configure it to use our MySQL database as the session store, and set a
	// lifetime of 12 hours (so that sessions automatically expire 12 hours
	// after first being created).
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		logger: logger,
		//we init a models.snippetmodel instance with the connection pool and add it to the application depencies
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  cache,
		sessionManager: sessionManager,
	}
	//now that we have a handler above (home) we need a router, in go termiology its called servemux
	//we will stop using the http.ListenAndServe, and we will use the http.Server struct for better control over our server
	//err = http.ListenAndServe(*addr, app.routes())
	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
		//we create a pointer from our structured logger handler which writes log entries at error level and assign it to the errorlog
		//field of the server.
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	//the value returned from the flag.String is a pointer to the flag value, not the value itself.
	//so we need to defrence it with *
	logger.Info("starting server", "addr", srv.Addr)
	err = srv.ListenAndServe()
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
