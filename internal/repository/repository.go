package repository

import "github.com/radu0v/1quoteEveryDay/internal/models"

type DataBaseRepo interface {
	AddQuote(quote string, author string) error
	GetQuotes() ([]models.Quote, error)
	SetDailyQuote() error
	IsSubscribed(emailAddr string) (bool, error)
	AddSubscriber(email string, name string) error
	GetSubscribers() ([]models.EmailData, error)
}
