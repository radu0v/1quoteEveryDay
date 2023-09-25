package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxLifetimeDB = 5 * time.Minute
const maxOpenDBConn = 10
const maxIdleDBConn = 5

func ConnectSql(dsn string) (*DB, error) {
	db, err := newDatabase(dsn)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(maxLifetimeDB)
	db.SetMaxIdleConns(maxIdleDBConn)
	db.SetMaxOpenConns(maxOpenDBConn)

	dbConn.SQL = db
	err = testDB(db)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

func testDB(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func newDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
