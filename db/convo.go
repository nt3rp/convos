package db

import (
	"database/sql"
	"strconv"

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
	Read      bool     `json:"read"`
	Children  []*Convo `json:"replies"`
}

func (c *Convo) Validate() bool {
	return true
}

func GetConvos(userId string) ([]*Convo, error) {
	db, err := DB()
	if err != nil {
		return nil, errgo.WithCausef(err, ErrConnection, "Error retrieving DB Connection")
	}

	rows, err := db.Query(`
		SELECT
		c.id, c.parent_id, c.sender_id, c.recipient_id, c.subject, c.body, r.user_id is not null
		FROM convos AS c
		LEFT JOIN read_status AS r ON r.thread_id = c.id AND r.user_id = $1
		WHERE c.parent_id = c.id
		AND (c.sender_id = $1 OR c.recipient_id = $1)
	`, userId)
	defer rows.Close()

	var cs []*Convo
	for rows.Next() {
		c := &Convo{}
		if err := rows.Scan(&c.Id, &c.Parent, &c.Sender, &c.Recipient, &c.Subject, &c.Body, &c.Read); err != nil {
			return cs, errgo.WithCausef(err, ErrRowScan, "Error Scanning Row")
		}

		cs = append(cs, c)
	}

	if err := rows.Err(); err != nil {
		return cs, errgo.WithCausef(err, ErrRowUnknown, "Unknown problem with `rows` object")
	}

	return cs, nil
}

func GetConvo(userId, convoId string) (*Convo, error) {
	db, err := DB()
	if err != nil {
		return nil, errgo.WithCausef(err, ErrConnection, "Error retrieving DB Connection")
	}

	c := &Convo{}
	err = db.QueryRow(`
		SELECT c.id, c.parent_id, c.sender_id, c.recipient_id, c.subject, c.body, r.user_id is not null
		FROM convos
		LEFT JOIN read_status AS r ON r.thread_id = c.id AND r.user_id = $2
		WHERE id = $1
		AND (c.sender_id = $2 OR c.recipient_id = $2)
	`, convoId, userId).Scan(
		&c.Id, &c.Parent, &c.Sender, &c.Recipient, &c.Subject, &c.Body, &c.Read,
	)

	if err == sql.ErrNoRows {
		return c, errgo.WithCausef(err, ErrNoRows, "Unable to find convo with id '%s'.", convoId)
	}

	if err != nil {
		return c, errgo.WithCausef(err, ErrRowScan, "Error Scanning Row")
	}

	return c, nil
}

func DeleteConvo(userId, convoId string) error {
	db, err := DB()
	if err != nil {
		return err
	}

	// No need to update read status on delete, should be handled by DB

	result, err := db.Exec(`
		DELETE
		FROM convos
		WHERE id = $1
		AND (sender_id = $2 OR recipient_id = $2)
	`, convoId, userId)

	if err != nil {
		return errgo.WithCausef(err, ErrRowDelete, "Error deleting convo.")
	}

	count, _ := result.RowsAffected()

	if err == sql.ErrNoRows || count == 0 {
		return errgo.WithCausef(err, ErrNoRows, "Unable to find convo with id '%s'.", convoId)
	}

	return nil
}

func CreateConvo(userId string, convo *Convo) (*Convo, error) {
	db, err := DB()
	if err != nil {
		return nil, errgo.WithCausef(err, ErrConnection, "Error retrieving DB Connection")
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, errgo.WithCausef(err, ErrTransaction, "Error starting transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	c := &Convo{}
	err = tx.QueryRow(`
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
		return nil, errgo.WithCausef(err, ErrRowCreate, "Error creating conversation")
	}

	result, err := tx.Exec(`
		INSERT INTO
		read_status (user_id, thread_id)
		VALUES ($1, $2)
	`, userId, c.Id)

	if err != nil {
		return nil, errgo.WithCausef(err, ErrRowCreate, "Error updating read status")
	}

	count, _ := result.RowsAffected()

	if err == sql.ErrNoRows || count == 0 {
		return nil, errgo.WithCausef(err, ErrRowCreate, "Unable to update read status")
	}

	c.Read = true

	return c, nil
}

func UpdateConvo(userId, convoId string, patch map[string]string) (*Convo, error) {
	db, err := DB()
	if err != nil {
		return nil, errgo.WithCausef(err, ErrConnection, "Error retrieving DB Connection")
	}

	val, ok := patch["read"]
	read, _ := strconv.ParseBool(val)
	if ok {
		var stmt string
		if read {
			stmt = "INSERT INTO read_status (user_id, thread_id) VALUES ($1, $2)"
		} else {
			stmt = "DELETE FROM read_status WHERE user_id = $1, thread_id = $2"
		}

		_, err := db.Exec(stmt, userId, convoId)

		if err != nil {
			return nil, errgo.WithCausef(err, ErrRowCreate, "Error updating read status")
		}
	}

	return GetConvo(userId, convoId)
}
