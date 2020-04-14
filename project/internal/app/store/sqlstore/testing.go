package sqlstore

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"testing"
)

func TestDB(t *testing.T, databaseURL string) (*sqlx.DB, func(...string)) {
	t.Helper()

	db, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	return db, func(tables ...string) {
		if len(tables) > 0 {
			_, err := db.Exec(fmt.Sprintf(`TRUNCATE bank.%s CASCADE`, strings.Join(tables, ", ")))
			if err != nil {
				t.Fatal(err)
			}
		}

		db.Close()
	}
}
