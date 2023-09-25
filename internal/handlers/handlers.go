package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/radu0v/1quoteEveryDay/internal/driver"
	"github.com/radu0v/1quoteEveryDay/internal/models"
	"github.com/radu0v/1quoteEveryDay/internal/render"
	"github.com/radu0v/1quoteEveryDay/internal/repository"
	"github.com/radu0v/1quoteEveryDay/internal/repository/postgres"
	"gopkg.in/gomail.v2"
)

// repository for the database functions
type Repository struct {
	DB repository.DataBaseRepo
}

func NewRepo(dbConn *driver.DB) *Repository {
	return &Repository{
		DB: postgres.NewPostgresDB(dbConn.SQL),
	}
}

var Repo *Repository

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	err := m.DB.SetDailyQuote()
	if err != nil {
		log.Println(err)
	}
	render.RenderTemplate(w, r, "home.page.tmpl", &models.Data{
		Quote:  models.DQ.Quote,
		Author: models.DQ.Author,
	})
}

func (m *Repository) PostHome(w http.ResponseWriter, r *http.Request) {
	models.Subscriber.Email = r.FormValue("email")
	models.Subscriber.Name = r.FormValue("name")

	// check if the email is already present in the database

	isSubscribed, err := m.DB.IsSubscribed(models.Subscriber.Email)
	if err != nil {
		log.Println("Error checking subscribers: ", err)
	}
	// if the user is not a subsriber , add it to the subscribers db table
	// and send greeting email
	if !isSubscribed {
		err := m.DB.AddSubscriber(models.Subscriber.Email, models.Subscriber.Name)
		if err != nil {
			log.Println("error adding subscriber to db table: ", err)
		}
		//send greeting email
		m := gomail.NewMessage()
		m.SetHeader("From", "aradu96.v@gmail.com")
		m.SetHeader("To", models.Subscriber.Email)

		m.SetHeader("Subject", "1 quote every day: Subscription")
		content := fmt.Sprintf("Hey %s! You are now subscribed to 1qed.com. You are set to receive one quote every day!", models.Subscriber.Name)
		m.SetBody("text/plain", content)

		d := gomail.NewDialer("smtp.gmail.com", 587, "aradu96.v@gmail.com", "dnpx pedw zfnr syuf")

		// Send the email to Bob, Cora and Dan.
		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}

	}
	render.RenderTemplate(w, r, "thank-you.page.tmpl", &models.Data{
		IsSubscribed: isSubscribed,
	})
}

// admin pages handlers

// dashboard for admin page
func (m *Repository) Admin(w http.ResponseWriter, r *http.Request) {
	err := m.DB.SetDailyQuote()
	if err != nil {
		log.Println(err)
	}
	render.RenderTemplate(w, r, "admin.page.tmpl", &models.Data{
		Quote:  models.DQ.Quote,
		Author: models.DQ.Author,
	})
}

// handler for admin page /quotes
func (m *Repository) AdminQuotes(w http.ResponseWriter, r *http.Request) {
	quotes, err := m.DB.GetQuotes()
	if err != nil {
		log.Println("Could not get quotes from database: ", err)
	}
	render.RenderTemplate(w, r, "admin-quotes.page.tmpl", &models.Data{
		Quotes: quotes,
	})
}

// handler for admin page /subscribers
func (m *Repository) Subscribers(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "admin.page.tmpl", &models.Data{})
}
