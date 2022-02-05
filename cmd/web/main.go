package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/tomdim/bookings/pkg/config"
	"github.com/tomdim/bookings/pkg/handlers"
	"github.com/tomdim/bookings/pkg/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8889"

var app config.AppConfig
var session *scs.SessionManager

// main is the main entrypoint
func main() {
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

	fmt.Printf("Starting application on port %s", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
