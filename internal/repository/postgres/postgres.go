package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/radu0v/1quoteEveryDay/internal/models"
	"github.com/radu0v/1quoteEveryDay/internal/repository"
)

type postgresDB struct {
	DB *sql.DB
}

func NewPostgresDB(conn *sql.DB) repository.DataBaseRepo {
	return &postgresDB{
		DB: conn,
	}
}

func (m *postgresDB) AddQuote(quote string, author string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into quotes(quote,author) values($1,$2)`
	_, err := m.DB.ExecContext(ctx, query, quote, author)
	if err != nil {
		log.Println("Error adding quote to database:", err)
		return err
	}
	return nil
}

func (m *postgresDB) DeleteQuote(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from quotes where id=$1`
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDB) GetQuotes() ([]models.Quote, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// get all the quotes in the db
	rows, err := m.DB.QueryContext(ctx, "select id, quote,author from quotes")
	if err != nil {
		log.Println("Error in getting quotes from database: ", err)
		return nil, err
	}
	defer rows.Close()
	var quote models.Quote
	var quotes []models.Quote
	for rows.Next() {
		err := rows.Scan(&quote.ID, &quote.Quote, &quote.Author)
		if err != nil {
			log.Println("Error scanning rows from quotes table in database: ", err)
			return []models.Quote{}, err
		}
		quotes = append(quotes, quote)
	}
	if err = rows.Err(); err != nil {
		return quotes, err
	}
	return quotes, nil
}

func (m *postgresDB) SetDailyQuote() error {
	// the daily quote is stored in the struct models.DQ
	// so if the date is the same as the current date
	// i won't do anything. If the date is different we
	// to set the daily quote in the struct

	//variables for all the function
	currentDate := time.Now().Format("2006-01-02")
	if models.DQ.Date != currentDate {
		// we check if the the quote_of_day table in db is empty
		isEmpty, err := m.tableIsEmpty("quote_of_day")
		if err != nil {
			log.Println("Error checking if table quote_of_day is empty: ", err)
			return err
		}
		if isEmpty {
			// if the db table is empty we have to insert the data in it
			// so we have to get all the quotes,choose one randomly and
			// insert it in the db table
			quotes, err := m.GetQuotes()
			if err != nil {
				log.Println("Error getting the quotes in SetDailyQuote() :", err)
				return err
			}
			//randomize the choice
			seed := rand.NewSource(time.Now().Unix())
			r := rand.New(seed)
			quote := quotes[r.Intn(len(quotes))]
			models.DQ.Quote = quote.Quote
			models.DQ.Author = quote.Author
			models.DQ.Date = currentDate
			//id := quote.ID

			// insert into table quote_of_day
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			query := `insert into quote_of_day(id_quote,date) values($1,$2)`
			_, err = m.DB.ExecContext(ctx, query, quote.ID, currentDate)
			if err != nil {
				log.Println("Error inserting new values in table quote_of_day.")
				return err
			}
		} else {
			id, date := m.getQOD()
			//case where the quote in the db table is valid
			//but is not present in the struct
			if date == currentDate {
				// we have to retrieve the quote by the id
				quote, author, err := m.getQuoteByID(id)
				if err != nil {
					log.Println("Error while retrieving quote by id in SetDailyQuote.")
					return err
				}
				models.DQ.Quote = quote
				models.DQ.Author = author
				models.DQ.Date = currentDate
			} else {
				//case where the quote in the db is not valid since it
				// does not have the current date;we have to update it.
				// get all the quotes and exclude the one in the table
				quotes, err := m.GetQuotes()
				if err != nil {
					log.Println("Error getting the quotes in SetDailyQuote() :", err)
					return err
				}
				//for loop where we look for a random quote that is not
				// the current one
				isSame := true
				var quote models.Quote
				for isSame {
					seed := rand.NewSource(time.Now().Unix())
					r := rand.New(seed)
					quote = quotes[r.Intn(len(quotes))]
					if quote.ID != id {
						isSame = false
					}
				}
				models.DQ.Quote = quote.Quote
				models.DQ.Author = quote.Author
				models.DQ.Date = currentDate

				//update the table
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				query := "update quote_of_day set id_quote=$1,date=$2 where id_quote=$3"
				_, err = m.DB.ExecContext(ctx, query, quote.ID, currentDate, id)
				if err != nil {
					log.Println("Error while trying to update quote_of_day table.")
					return err
				}
			}
		}
	}
	return nil
}

func (m *postgresDB) getQOD() (id int, date string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var quote_id int
	var date_db string

	rows, err := m.DB.QueryContext(ctx, `select id_quote,date from quote_of_day`)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&quote_id, &date_db)
		if err != nil {
			log.Println(err)
		}
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
	}
	return quote_id, date_db
}

// checks whether the table quote_of_day is empty
func (m *postgresDB) tableIsEmpty(table string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var result int
	query := fmt.Sprintf("select count(*) from %s", table)
	row := m.DB.QueryRowContext(ctx, query)

	err := row.Scan(&result)
	if err != nil {
		return false, err
	}
	if result == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (m *postgresDB) getQuoteByID(id int) (quote string, author string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select quote,author from quotes where id=$1`
	row := m.DB.QueryRowContext(ctx, query, id)

	err = row.Scan(&quote, &author)
	if err != nil {
		log.Println("error getting quote by id:", err)
		return
	}
	return
}

// function that checks whether an email is already present in the db
func (m *postgresDB) IsSubscribed(emailAddr string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var isSubsribed bool

	// check if the subscribers db table is empty or not
	// if it is empty return false
	tableEmpty, err := m.tableIsEmpty("subscribers")
	if err != nil {
		return isSubsribed, err
	}
	if tableEmpty {
		return false, nil
	}

	// get all the email and name in the db and save them in a map
	// then check whether the given email exist or not
	rows, err := m.DB.QueryContext(ctx, "select email, name from subscribers")
	if err != nil {
		return isSubsribed, err
	}
	defer rows.Close()
	var email, name string
	subscribers := map[string]string{}
	for rows.Next() {
		err := rows.Scan(&email, &name)
		if err != nil {
			return isSubsribed, err
		}
		subscribers[email] = name
	}
	if err = rows.Err(); err != nil {
		return isSubsribed, err
	}

	_, isSubsribed = subscribers[emailAddr]

	return isSubsribed, nil
}

func (m *postgresDB) AddSubscriber(email string, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into subscribers (email,name,date) values($1,$2,$3)`

	_, err := m.DB.ExecContext(ctx, query, email, name, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDB) GetSubscribers() ([]models.EmailData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	subs := []models.EmailData{}
	rows, err := m.DB.QueryContext(ctx, "select email,name from subscribers")
	if err != nil {
		return subs, err
	}
	defer rows.Close()
	var email, name string
	for rows.Next() {
		err := rows.Scan(&email, &name)
		if err != nil {
			return subs, err
		}
		subs = append(subs, models.EmailData{Email: email, Name: name})
	}
	if err := rows.Err(); err != nil {
		return subs, err
	}
	return subs, nil
}

func (m *postgresDB) Unsubscribe(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from subscribers * where email=$1`
	_, err := m.DB.ExecContext(ctx, query, email)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDB) Authenticate(user string, pass string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// check if the user exists in the db table
	var username, password string
	row := m.DB.QueryRowContext(ctx, "select user,password from users where id=1")
	err := row.Scan(&username, &password)
	if err != nil {
		return err
	}
	if err = row.Err(); err != nil {
		return err
	}
	if username == user {
		match, err := argon2id.ComparePasswordAndHash(pass, password)
		if err != nil {
			return err
		}
		if !match {
			return errors.New("invalid credentials")
		}
		return nil
	} else {
		return errors.New("invalid credentials")
	}
}
