package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID               int64
	Email            string
	EncrypedPassword string
	Name             string
	CreatedAt        time.Time
}

func GetUser(db *sql.DB, id int64) (*User, error) {
	u := &User{}
	err := db.QueryRow(
		"SELECT id, email, encryped_password, name FROM users WHERE id = ? LIMIT 1",
		id,
	).
		Scan(
			&u.ID,
			&u.Email,
			&u.EncrypedPassword,
			&u.Name,
		)
	if err != nil {
		return nil, err
	}
	return u, nil
}
