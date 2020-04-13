package model

type User struct {
	ID                int
	Login             string
	Name              string
	Surname           string
	Patronymic        string
	Passport          string
	EncryptedPassword string
}
