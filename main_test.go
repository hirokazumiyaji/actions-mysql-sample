package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hirokazumiyaji/actions-mysql-sample/database"
	"github.com/schemalex/schemalex"
	"github.com/schemalex/schemalex/diff"
)

func TestMain(m *testing.M) {
	conn, err := database.Open()
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		log.Fatalf("failed to ping: %v", err)
	}

	type query struct {
		stmt string
		args []interface{}
	}

	from, err := schemalex.NewSchemaSource(
		fmt.Sprintf(
			"mysql://%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_PASSWORD"),
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_NAME"),
		),
	)
	if err != nil {
		log.Fatalf("failed to new schema source: %v", err)
	}
	to, err := schemalex.NewSchemaSource("./docker/mysql/sql/initialize.sql")
	if err != nil {
		log.Fatalf("failed to new schema source: %v", err)
	}

	stmts := &bytes.Buffer{}
	p := schemalex.New()
	err = diff.Sources(stmts, from, to, diff.WithTransaction(true), diff.WithParser(p))
	if err != nil {
		log.Fatalf("failed to diff: %v", err)
	}

	qs := make([]*query, 0)
	for _, stmt := range strings.Split(stms, ";") {
		stmt = strings.TrimSpace(stmt)
		if len(stmt) == 0 {
			continue
		}
		qs = append(
			qs,
			&query{
				stmt: stmt[0],
				args: stmt[1:],
			},
		)
	}

	tx := conn.Begin()
	for _, q := range qs {
		_, err := tx.Exec(q.stmt, q.args...)
		if err != nil {
			log.Printf("failed to exec sql: %v", err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()

	os.Exit(m.Run())
}
