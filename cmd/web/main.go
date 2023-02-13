package main

import (
	"encoding/gob"
	"fmt"
	"github/Atul-Ranjan12/booking/internal/config"
	"github/Atul-Ranjan12/booking/internal/handlers"
	"github/Atul-Ranjan12/booking/internal/models"
	"github/Atul-Ranjan12/booking/internal/render"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":4000"
var app config.AppConfig
var session *scs.SessionManager

// main is the main function
func main() {

	err := run();
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting server on local host ");

	server := &http.Server {
		Addr: portNumber,
		Handler: routes(&app),
	}

	serverError := server.ListenAndServe()
	log.Fatal(serverError)
}

// Funciton to run the main application
func run() error {

	// Primitive types for the Session
	gob.Register(models.Reservation{})

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)
	
	return nil
}