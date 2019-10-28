package main

import (
	"log"
	"os"
	"testing"

	"github.com/hirokazumiyaji/actions-mysql-sample/database"
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

	os.Exit(m.Run())
}
