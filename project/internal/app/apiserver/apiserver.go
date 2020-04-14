package apiserver

import (
	"github.com/ENSLERMAN/soft-eng/project/internal/app/store/sqlstore"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/http"
)

func Start(config *Config) error {

	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	store := sqlstore.New(db)
	srv := newServer(store)

	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	defer db.Close()
	return db, err
}
