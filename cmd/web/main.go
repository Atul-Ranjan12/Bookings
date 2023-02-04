package main

import (
	"fmt"
	"github/Atul-Ranjan12/booking/pkg/config"
	"github/Atul-Ranjan12/booking/pkg/handlers"
	"github/Atul-Ranjan12/booking/pkg/render"
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
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	server := &http.Server {
		Addr: portNumber,
		Handler: routes(&app),
	}

	fmt.Println("Starting server on local host ");
	serverError := server.ListenAndServe()
	log.Fatal(serverError)
}