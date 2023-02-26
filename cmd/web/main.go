package main

import (
	"encoding/gob"
	"fmt"
	"github/Atul-Ranjan12/booking/internal/config"
	"github/Atul-Ranjan12/booking/internal/driver"
	"github/Atul-Ranjan12/booking/internal/handlers"
	"github/Atul-Ranjan12/booking/internal/helpers"
	"github/Atul-Ranjan12/booking/internal/models"
	"github/Atul-Ranjan12/booking/internal/render"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":4000"

var app config.AppConfig
var session *scs.SessionManager

// Initializing the logs
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {

	db, err := run()
	if err != nil {
		log.Fatal("AN unexpected error happened while connecting to the databse, quitting... with ", err)
	}
	defer db.SQL.Close()
	// Closing the Mailing channel
	defer close(app.MailChan)
	fmt.Println("Starting Mail Channel..")
	listenForMail()

	// msg := models.MailData{
	// 	To:      "john@do.ca",
	// 	From:    "me@here.com",
	// 	Subject: "Some subject",
	// 	Content: "",
	// }

	// app.MailChan <- msg

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting server on local host ")

	server := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	serverError := server.ListenAndServe()
	log.Fatal(serverError)
}

// Funciton to run the main application
func run() (*driver.DB, error) {

	// Primitive types for the Session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// Make a mailing channel
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=atulranjan password=")
	if err != nil {
		log.Fatal("Cannot connect to the database, dyring..")
	}
	log.Println("Connected to the database")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelper(&app)

	return db, nil
}
