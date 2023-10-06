package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/radu0v/1quoteEveryDay/internal/config"
	"github.com/radu0v/1quoteEveryDay/internal/driver"
	"github.com/radu0v/1quoteEveryDay/internal/handlers"
	"github.com/radu0v/1quoteEveryDay/internal/render"
	"github.com/robfig/cron/v3"
)

const portNumber = ":8080"

var app config.AppConfig
var sessionManager *scs.SessionManager

func main() {
	app.InProduction = false
	//initialize session and configure the session lifetime and cookie
	sessionManager = scs.New()
	sessionManager.Lifetime = 1 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = app.InProduction
	sessionManager.Cookie.HttpOnly = true

	app.Session = sessionManager
	// ask for database info
	fmt.Println("Hello!")
	//fmt.Println("Database host:")
	var dbHost = "localhost"
	//fmt.Scan(&dbHost)
	//fmt.Println("Database port:")
	var dbPort = "5432"
	//fmt.Scan(&dbPort)
	//fmt.Println("Database name:")
	var dbName = "1quote"
	//fmt.Scan(&dbName)
	//fmt.Println("Database user:")
	var dbUser = "radu"
	//fmt.Scan(&dbUser)
	//fmt.Println("Database password:")
	var dbPassword = ""
	//fmt.Scan(&dbPassword)
	// connect to database
	fmt.Println("Connecting to database...")
	db, err := driver.ConnectSql("host= " + dbHost + " port=" + dbPort + " dbname=" + dbName + " user=" + dbUser + " password=" + dbPassword)
	if err != nil {
		log.Fatal("Cannot connect to database!")
	}
	defer db.SQL.Close()
	log.Println("Connected!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache: ", err)
	}
	app.TemplateCache = tc
	app.UseCache = true
	repo := handlers.NewRepo(db, &app)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)

	// cronjob for sending a quote everyday
	c := cron.New()
	c.AddFunc("0 9 * * *", func() { sendNewsLetter() })
	go c.Start()
	c.Stop()

	// server
	fmt.Println("Starting application on port 8080")
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
