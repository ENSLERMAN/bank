package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ENSLERMAN/soft-eng/project/internal/app/model"
	"net/http"
	"strconv"
)

// handleUsersCreate обработчик - создание клиента.
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
		// записываем данные in json
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

		// кидаем запрос на создание клиента, если не получилось скидываем ошибку 422.
		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		// Sanitaze - нужен для того, шоб затереть пароль после создания клиента.
		u.Sanitaze()
		s.respond(w, r, http.StatusCreated, u)
	}
}

// handleBillCreate обработчик - создание счета пользователя.
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

		// берем из куки ид юзера.
		session, err := s.sessionStore.Get(r, "bank-system")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		userID, _ := strconv.Atoi(fmt.Sprint(session.Values["user_id"]))

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

// handleSessionsCreate обработчик - создание сессии, проще говоря авторизация юзера.
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

// authenticateUser - обработчик сессии, проверка ид юзера с бд.
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

// handleBillDelete - обработчик закрытия счета.
func (s *server) handleBillDelete() http.HandlerFunc {

	type request struct {
		BillID int `json:"bill_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		// получаем ид юзера из куки.
		session, err := s.sessionStore.Get(r, "bank-system")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		userID, _ := strconv.Atoi(fmt.Sprint(session.Values["user_id"]))

		// если не подходит, кидаем ошибку 401.
		_, err = s.store.Bill().FindByUser(userID, req.BillID)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		// если все подошло, закрываем счет клиента ( просто удаляем ).
		if err := s.store.Bill().DeleteBill(req.BillID); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

// handleGetAllUserBills - обработчик, получаем все счета клиента.
func (s *server) handleGetAllUserBills() http.HandlerFunc {


	return func(w http.ResponseWriter, r *http.Request) {

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

func (s *server) handleGetAllUserPayments() http.HandlerFunc {

	type request struct {
		Number     int  `json:"number"`
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

		_, err = s.store.Bill().FindByBill(userID, req.Number)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.store.Bill().GetAllUserPayments(req.Number)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, u)
	}
}

// handleSendMoney - обработчик перевода денег.
func (s *server) handleSendMoney() http.HandlerFunc {

	type request struct {
		BillID     int  `json:"bill_id"`
		NumberDest int  `json:"number_dest"`
		Amount     uint `json:"amount"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		// получаем ид юзера из сессии.
		session, err := s.sessionStore.Get(r, "bank-system")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		userID, _ := strconv.Atoi(fmt.Sprint(session.Values["user_id"]))

		_, err = s.store.Bill().FindByUser(userID, req.BillID)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		if err := s.store.Bill().TransferMoney(req.NumberDest, req.Amount, req.BillID); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}
