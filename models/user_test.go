package models

import (
	"database/sql"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hirokazumiyaji/actions-mysql-sample/database"
	"golang.org/x/crypto/bcrypt"
)

func TestGetUser(t *testing.T) {
	conn, err := database.Open()
	if err != nil {
		t.Errorf("failed to open database: %v", err)
		return
	}
	defer conn.Close()

	actual, err := GetUser(conn, 0)
	if sql.ErrNoRows != err {
		t.Errorf("expected: %v, actual: %v", sql.ErrNoRows, err)
		return
	}
	if actual != nil {
		t.Errorf("expected: nil, actual: %v", actual)
		return
	}

	encrypedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		t.Errorf("failed to encrypt password: %v", err)
		return
	}
	expected := &User{
		Email:            "admin@example.com",
		EncrypedPassword: string(encrypedPassword),
		Name:             "admin",
	}
	stmt, err := conn.Prepare("INSERT INTO users (email, encryped_password, name) VALUES (?, ?, ?)")
	if err != nil {
		t.Errorf("failed to prepare: %v", err)
		return
	}
	row, err := stmt.Exec(expected.Email, expected.EncrypedPassword, expected.Name)
	if err != nil {
		stmt.Close()
		t.Errorf("failed to exec sql: %v", err)
		return
	}
	expected.ID, err = row.LastInsertId()
	if err != nil {
		stmt.Close()
		t.Errorf("failed to last insert id: %v", err)
		return
	}
	stmt.Close()

	actual, err = GetUser(conn, expected.ID)
	if err != nil {
		t.Errorf("failed to get user: %v", err)
		return
	}

	if !cmp.Equal(expected, actual) {
		t.Error(cmp.Diff(expected, actual))
	}
}
