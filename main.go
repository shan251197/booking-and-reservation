package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/shan251197/booking-and-reservation/config"
	"github.com/shan251197/booking-and-reservation/handlers"
	"github.com/shan251197/booking-and-reservation/render"
)

const portNumber = ":8000"

var app config.AppConfig

var session *scs.SessionManager

func main() {

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Secure = app.InProduction
	session.Cookie.SameSite = http.SameSiteLaxMode

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc

	app.UseCache = false

	repo := handlers.NewRepo(&app)

	handlers.NewHandlers(repo)

	render.NewTemplate(&app)

	fmt.Printf("%s", fmt.Sprintf("starting application on port %s \n", portNumber))

	// http.ListenAndServe(":8000", nil)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()

	log.Fatal(err)
}
