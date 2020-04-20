package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ENSLERMAN/soft-eng/project/internal/app/model"
	"net/http"
	"strconv"
)

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Login      string `json:"login"`
		Name       string `json:"name"`
		Surname    string `json:"surname"`
		Patronymic string `json:"patronymic"`
		Passport   string `json:"passport"`
		Password   string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Login:      req.Login,
			Name:       req.Name,
			Surname:    req.Surname,
			Patronymic: req.Patronymic,
			Passport:   req.Passport,
			Password:   req.Password,
		}

		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitaze()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleBillCreate() http.HandlerFunc {
	type request struct {
		Type int `json:"type"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		session, err := s.sessionStore.Get(r, "bank-system")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		userID, _ := strconv.Atoi(fmt.Sprint(session.Values["user_id"]))
		fmt.Println(userID)

		u := &model.Bill{
			Type:   req.Type,
			Number: model.RandomizeCardNumber(),
		}

		if err := s.store.Bill().CreateBill(u, userID); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByLogin(req.Login)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectLoginOrPassword)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.store.User().FindByID(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) handleBillDelete() http.HandlerFunc {
	type request struct {
		UserID int `json:"user_id"`
		BillID int `json:"bill_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Bill().DeleteBill(req.BillID); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleGetAllUserBills() http.HandlerFunc {
	type request struct {
		UserID string `json:"user_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		session, err := s.sessionStore.Get(r, "bank-system")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		userID, _ := strconv.Atoi(fmt.Sprint(session.Values["user_id"]))

		u, err := s.store.Bill().GetAllUserBills(userID)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, u)
	}
}
