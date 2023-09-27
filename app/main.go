package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/radu0v/1quoteEveryDay/internal/driver"
	"github.com/radu0v/1quoteEveryDay/internal/handlers"
	"github.com/robfig/cron/v3"
)

const portNumber = ":8080"

func main() {
	// ask for database info
	fmt.Println("Hello!")
	//fmt.Println("Database host:")
	var dbHost = "localhost"
	//fmt.Scan(&dbHost)
	//fmt.Println("Database port:")
	var dbPort = ""
	//fmt.Scan(&dbPort)
	//fmt.Println("Database name:")
	var dbName = ""
	//fmt.Scan(&dbName)
	//fmt.Println("Database user:")
	var dbUser = ""
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

	repo := handlers.NewRepo(db)
	handlers.NewHandlers(repo)

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
