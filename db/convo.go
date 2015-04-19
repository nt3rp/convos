package db

import (
	_ "github.com/lib/pq"
)

type Convo struct {
	Id        int      `json:"id"`
	Sender    int      `json:"sender"`
	Recipient int      `json:"recipient"`
	Subject   string   `json:"subject"`
	Body      string   `json:"body"`
	Status    string   `json:"status"`
	Children  []*Convo `json:"replies"`
}

func GetConvos() ([]*Convo, error){
	db, err := DB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
		SELECT id, sender_id, recipient_id, subject, body
		FROM convos
		WHERE parent_id = id
	`)
	defer rows.Close()

	var cs []*Convo
	for rows.Next() {
		c := &Convo{}
		if err := rows.Scan(&c.Id, &c.Sender, &c.Recipient, &c.Subject, &c.Body); err != nil {
			return cs, err
		}

		cs = append(cs, c)
	}

	if err := rows.Err(); err != nil {
		return cs, err
	}

	return cs, err
}

func GetConvo(id string) (*Convo, error) {
	db, err := DB()
	if err != nil {
		return nil, err
	}

	c := &Convo{}
	err = db.QueryRow(`
		SELECT id, sender_id, recipient_id, subject, body
		FROM convos
		WHERE id = $1
	`, id).Scan(
		&c.Id, &c.Sender, &c.Recipient, &c.Subject, &c.Body,
	)

	return c, err
}

func DeleteConvo(id string) error {
	db, err := DB()
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		DELETE
		FROM convos
		WHERE id = $1
	`, id)

	// TODO: Check rows affected

	return err
}

func CreateConvo(convo *Convo) error {
	db, err := DB()
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO
		convos (parent_id, sender_id, recipient_id, subject, body)
		VALUES (lastval(), $1, $2, $3, $4)
	`, convo.Sender, convo.Recipient, convo.Subject, convo.Body)

	// TODO: Check rows affected

	return err
}
