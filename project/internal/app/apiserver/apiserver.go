// Пакет apiserser позволяет конфигурировать наш сервер и запускать его.
package apiserver

import (
	"github.com/ENSLERMAN/soft-eng/project/internal/app/store/sqlstore"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Start запускаем сервер.
func Start(config *Config) error {

	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)
	sessionsStore := sessions.NewCookieStore([]byte(config.SessionKey))
	srv := newServer(store, sessionsStore)
	logrus.Infof("started server http://localhost%s", config.BindAddr)

	return http.ListenAndServe(config.BindAddr, srv)
}

// newDB открываем соединение с бд
func newDB(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
