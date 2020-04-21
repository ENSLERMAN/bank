package model

import "testing"

// TestUser - создаем нового юзера и сразу его отдаем
func TestUser(t *testing.T) *User {
	return &User{
		Login:      "anakonda3000",
		Password:   "sobakapes123",
		Name:       "Ivan",
		Surname:    "Ivanov",
		Patronymic: "Ivanovich",
		Passport:   "0001",
	}
}
