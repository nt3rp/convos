package db

import (
	"fmt"

	"github.com/juju/errgo"
	_ "github.com/lib/pq"
)

// These functions are not intended for public consumption but exposed so tests have access
func TruncateTable(name string) error {
	db, err := DB()
	if err != nil {
		return errgo.WithCausef(err, ErrConnection, "Error retrieving DB Connection")
	}

	stmt := fmt.Sprintf("DELETE FROM %s", name)
	_, err = db.Exec(stmt)

	if err != nil {
		return errgo.WithCausef(err, ErrTruncate, "Error truncating table")
	}

	return nil
}

func AddUser(userId, name string) error {
	db, err := DB()
	if err != nil {
		return errgo.WithCausef(err, ErrConnection, "Error retrieving DB Connection")
	}

	_, err = db.Exec(`
		INSERT INTO users (id, fullname) VALUES ($1, $2)
	`, userId, name)

	if err != nil {
		return errgo.WithCausef(err, ErrRowCreate, "Error adding user.")
	}

	return nil
}
