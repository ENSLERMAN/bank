package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://ensler:gavnojopa@localhost:5432/homedb?sslmode=disable"

	}

	os.Exit(m.Run())
}
