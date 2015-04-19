package db

import (
	"database/sql"
	"log"
)

var (
	dbSession *sql.DB
)

func DB() (*sql.DB, error) {
	if dbSession != nil {
		return dbSession, nil
	}

	// TODO: Make params configurable
	db, err := sql.Open("postgres", "dbname=convos sslmode=disable")
	if err != nil {
		log.Fatal(err)
	} else {
		dbSession = db
	}

	return db, err
}
