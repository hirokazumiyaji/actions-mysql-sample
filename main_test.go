package main

import (
	"log"
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

	os.Exist(m.Run())
}
