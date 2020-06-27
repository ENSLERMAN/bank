package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ENSLERMAN/soft-eng/project/internal/app/model"
	"github.com/ENSLERMAN/soft-eng/project/internal/app/store"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

const (
	sessionName        = "bank-system"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectLoginOrPassword = errors.New("incorrect login or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

// newServer - создаем новый сервер из параметров.
func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// configureRouter - здесь записываем нужные нам роуты.
func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"})
	originsOk := handlers.AllowedOrigins([]string{"http://bank.enslerman.ru","http://localhost:4200"})
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	cockSucker := handlers.AllowCredentials()
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST", "OPTIONS")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST", "OPTIONS")

	s.router.Use(handlers.CORS(methodsOk, originsOk, headersOk, cockSucker))

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/whoami", s.handleWhoami()).Methods("GET")
	private.HandleFunc("/create_bill", s.handleBillCreate()).Methods("POST", "OPTIONS")
	private.HandleFunc("/bills", s.handleGetAllUserBills()).Methods("GET", "OPTIONS")
	private.HandleFunc("/get_bills", s.GetRestOfTheBills()).Methods("GET", "OPTIONS")
	private.HandleFunc("/bills/{id}", s.GetUserBillByID()).Methods("GET", "OPTIONS")
	private.HandleFunc("/bills/{id}", s.handleBillDelete()).Methods("DELETE", "OPTIONS")
	private.HandleFunc("/send_money", s.handleSendMoney()).Methods("POST", "OPTIONS")
	private.HandleFunc("/payments", s.handleGetAllUserPayments()).Methods("GET", "OPTIONS")
	private.HandleFunc("/payments/{id}", s.GetUserPaymentsByID()).Methods("GET", "OPTIONS")
	private.HandleFunc("/get_money/{id}", s.GetMoney()).Methods("GET", "OPTIONS")
}

// handleWhoami - обработчик для проверки юзера.
func (s *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}

// logRequest - middleware для логирования действий.
func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

// setRequestID - middleware обработчик для создания каждому request уникального id.
func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

// error - банальный вывод ошибки.
func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

// respond - отдача ответа in json.
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
