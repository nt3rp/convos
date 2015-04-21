package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/juju/errgo"
)

var (
	dbSession *sql.DB
)

func Initialize(dbName string) {
	db, err := Connect(dbName)

	if err != nil {
		dbSession = db
	}
}

func Connect(dbName string) (*sql.DB, error) {
	connParams := fmt.Sprintf("dbname=%s sslmode=disable", dbName)
	db, err := sql.Open("postgres", connParams)
	if err != nil {
		log.Fatal(err)
	} else {
		dbSession = db
	}

	return db, err
}

func DB() (*sql.DB, error) {
	if dbSession != nil {
		return dbSession, nil
	}

	return nil, errgo.WithCausef(ErrUninitialized, ErrUninitialized, "Database was not initialized")
}
