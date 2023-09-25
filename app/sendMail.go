package main

import (
	"fmt"
	"log"

	"github.com/radu0v/1quoteEveryDay/internal/handlers"
	"github.com/radu0v/1quoteEveryDay/internal/models"
	"gopkg.in/gomail.v2"
)

func sendNewsLetter() {
	err := handlers.Repo.DB.SetDailyQuote()
	if err != nil {
		log.Println("Error setting daily quote:", err)
	}
	d := gomail.NewDialer("smtp.gmail.com", 587, "", "")
	s, err := d.Dial()
	if err != nil {
		panic(err)
	}
	m := gomail.NewMessage()

	subs, err := handlers.Repo.DB.GetSubscribers()
	if err != nil {
		log.Println("Error getting subscribers table: ", err)
	}

	for _, sub := range subs {
		m.SetHeader("From", "")
		m.SetAddressHeader("To", sub.Email, sub.Name)
		m.SetHeader("Subject", "Daily quote from 1 quote every day")
		m.SetBody("text/plain", fmt.Sprintf("%s\n%s", models.DQ.Quote, models.DQ.Author))
		if err := gomail.Send(s, m); err != nil {
			log.Printf("Could not send email to %s:%v", sub.Email, err)
		}
		m.Reset()
	}

}
