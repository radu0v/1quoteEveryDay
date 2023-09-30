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
		msg := gomail.NewMessage()
		msg.SetHeader("From", "")
		msg.SetHeader("To", models.Subscriber.Email)

		msg.SetHeader("Subject", "1 quote every day: Subscription")
		content := fmt.Sprintf("Hey %s! You are now subscribed. You are set to receive one quote every day!", models.Subscriber.Name)
		msg.SetBody("text/plain", content)

		d := gomail.NewDialer("smtp.gmail.com", 587, "", "")

		// Send the email
		if err := d.DialAndSend(msg); err != nil {
			panic(err)
		}

	}
	render.RenderTemplate(w, r, "post-subscription.page.tmpl", &models.Data{
		IsSubscribed: isSubscribed,
	})
}

func (m *Repository) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "unsubscribe.page.tmpl", &models.Data{})
}

func (m *Repository) UnsubscribePost(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	//if email is not present in database
	// write to the web page that the user will be unsubscribed and
	// they will receive an email confirming they are no longer subscribed
	unsubscribed := false
	isSubscribed, err := m.DB.IsSubscribed(email)
	if err != nil {
		log.Println("error checking if user is subscribed: ", err)
	}

	if isSubscribed {
		err := m.DB.Unsubscribe(email)
		if err != nil {
			log.Println("Error unsubscribing user: ", err)
		}
		unsubscribed = true
		//send email to confirm that user is no longer subscribed
		render.RenderTemplate(w, r, "unsubscribe.page.tmpl", &models.Data{
			Bool: unsubscribed,
		})
	} else {
		render.RenderTemplate(w, r, "unsubscribe.page.tmpl", &models.Data{
			Bool: unsubscribed,
		})
	}

}

func (m *Repository) Feedback(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "feedback.page.tmpl", &models.Data{})
}

func (m *Repository) FeedbackPost(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	name := r.FormValue("name")
	messagge := r.FormValue("message")

	//send mail with the feedback
	msg := gomail.NewMessage()
	msg.SetHeader("To", "")
	msg.SetHeader("From", "")
	msg.SetHeader("Subject", fmt.Sprintf("Feedback from %s (%s)", name, email))
	msg.SetBody("text/plain", messagge)

	d := gomail.NewDialer("smtp.gmail.com", 587, "", "")

	// send mail
	if err := d.DialAndSend(msg); err != nil {
		panic(err)
	}

	render.RenderTemplate(w, r, "feedback.page.tmpl", &models.Data{})

}

func (m *Repository) PrivacyPolicy(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "privacy-policy.page.tmpl", &models.Data{})
}

// admin pages handlers

// login page handler
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "login.page.tmpl", &models.Data{})
}

// post login page handler
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	//check user and password
	user := r.FormValue("username")
	pass := r.FormValue("password")

	err := m.DB.Authenticate(user, pass)
	if err != nil {
		log.Println("error authenticating user:", err)

	}

}

// dashboard for admin page
func (m *Repository) Admin(w http.ResponseWriter, r *http.Request) {
	err := m.DB.SetDailyQuote()
	if err != nil {
		log.Println(err)
	}

	quotes, err := m.DB.GetQuotes()
	if err != nil {
		log.Println("error getting quotes:", err)
	}
	nrQuotes := len(quotes)

	subs, err := m.DB.GetSubscribers()
	if err != nil {
		log.Println("error getting subscribers: ", err)
	}
	nrSub := len(subs)
	lastSub := subs[len(subs)-1]
	intMap := map[string]int{
		"nrQuotes": nrQuotes,
		"nrSubs":   nrSub,
	}
	stringMap := map[string]string{
		"email": lastSub.Email,
	}
	render.RenderTemplate(w, r, "admin.page.tmpl", &models.Data{
		Quote:     models.DQ.Quote,
		Author:    models.DQ.Author,
		IntMap:    intMap,
		StringMap: stringMap,
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
	subs, err := m.DB.GetSubscribers()
	if err != nil {
		fmt.Println("error getting subscribers: ", err)
	}

	render.RenderTemplate(w, r, "admin-subscribers.page.tmpl", &models.Data{
		Subscribers: subs,
	})
}
