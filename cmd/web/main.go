package main

import (
	"encoding/gob"
	"fmt"
	"github.com/tomdim/bookings/internal/config"
	"github.com/tomdim/bookings/internal/handlers"
	"github.com/tomdim/bookings/internal/models"
	"github.com/tomdim/bookings/internal/render"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8889"

var app config.AppConfig
var session *scs.SessionManager

// main is the main entrypoint
func main() {
	// what am I going to put in the session
	gob.Register(models.Reservation{})

	// Change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	fmt.Printf("Starting application on port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
