package model

import "testing"

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
