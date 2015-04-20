package db

import (
	"database/sql"

	"github.com/juju/errgo"
	_ "github.com/lib/pq"
)

type Convo struct {
	Id        int      `json:"id"`
	Sender    int      `json:"sender"`
	Recipient int      `json:"recipient"`
	Parent    int      `json:"parent"`
	Subject   string   `json:"subject"`
	Body      string   `json:"body"`
	Status    string   `json:"status"`
	Children  []*Convo `json:"replies"`
}

func (c *Convo) Validate() bool {
	return true
}

func GetConvos() ([]*Convo, error) {
	db, err := DB()
	if err != nil {
		return nil, errgo.WithCausef(err, ErrConnection, "Error retrieving DB Connection")
	}

	rows, err := db.Query(`
		SELECT id, parent_id, sender_id, recipient_id, subject, body
		FROM convos
		WHERE parent_id = id
	`)
	defer rows.Close()

	var cs []*Convo
	for rows.Next() {
		c := &Convo{}
		if err := rows.Scan(&c.Id, &c.Parent, &c.Sender, &c.Recipient, &c.Subject, &c.Body); err != nil {
			return cs, errgo.WithCausef(err, ErrRowScan, "Error Scanning Row")
		}

		cs = append(cs, c)
	}

	if err := rows.Err(); err != nil {
		return cs, errgo.WithCausef(err, ErrRowUnknown, "Unknown problem with `rows` object")
	}

	return cs, nil
}

func GetConvo(id string) (*Convo, error) {
	db, err := DB()
	if err != nil {
		return nil, errgo.WithCausef(err, ErrConnection, "Error retrieving DB Connection")
	}

	c := &Convo{}
	err = db.QueryRow(`
		SELECT id, parent_id, sender_id, recipient_id, subject, body
		FROM convos
		WHERE id = $1
	`, id).Scan(
		&c.Id, &c.Parent, &c.Sender, &c.Recipient, &c.Subject, &c.Body,
	)

	if err == sql.ErrNoRows {
		return c, errgo.WithCausef(err, ErrNoRows, "Unable to find convo with id '%s'.", id)
	}

	if err != nil {
		return c, errgo.WithCausef(err, ErrRowScan, "Error Scanning Row")
	}

	return c, nil
}

func DeleteConvo(id string) error {
	db, err := DB()
	if err != nil {
		return err
	}

	result, err := db.Exec(`
		DELETE
		FROM convos
		WHERE id = $1
	`, id)

	count, _ := result.RowsAffected()

	if err == sql.ErrNoRows || count == 0 {
		return errgo.WithCausef(err, ErrNoRows, "Unable to find convo with id '%s'.", id)
	}

	if err != nil {
		return errgo.WithCausef(err, ErrRowDelete, "Error deleting convo.")
	}

	return nil
}

func CreateConvo(convo *Convo) (*Convo, error) {
	db, err := DB()
	if err != nil {
		return nil, errgo.WithCausef(err, ErrConnection, "Error retrieving DB Connection")
	}

	c := &Convo{}
	err = db.QueryRow(`
		INSERT INTO
		convos (parent_id, sender_id, recipient_id, subject, body)
		VALUES (
			CASE
			    WHEN $1=0 THEN lastval()
				ELSE $1
			END,
			$2, $3, $4, $5
		)
		RETURNING id, parent_id, sender_id, recipient_id, subject, body
	`, convo.Parent, convo.Sender, convo.Recipient, convo.Subject, convo.Body).Scan(
		&c.Id, &c.Parent, &c.Sender, &c.Recipient, &c.Subject, &c.Body,
	)

	if err != nil {
		return nil, errgo.WithCausef(err, ErrRowCreate, "Error Creating Rows")
	}

	return c, nil
}

func UpdateConvo(id, body string) (*Convo, error) {
	db, err := DB()
	if err != nil {
		return nil, errgo.WithCausef(err, ErrConnection, "Error retrieving DB Connection")
	}

	c := &Convo{}
	err = db.QueryRow(`
		UPDATE convos
		SET body = $2
		WHERE id = $1
		RETURNING id, sender_id, recipient_id, subject, body
	`, id, body).Scan(
		&c.Id, &c.Sender, &c.Recipient, &c.Subject, &c.Body,
	)

	if err != nil {
		return nil, errgo.WithCausef(err, ErrRowUpdate, "Error Updating Row")
	}

	return c, nil
}
