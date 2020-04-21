package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

// User - структура юзера
type User struct {
	ID                int    `json:"id"`
	Login             string `json:"login"`
	Name              string `json:"name"`
	Surname           string `json:"surname"`
	Patronymic        string `json:"patronymic"`
	Passport          string `json:"passport"`
	EncryptedPassword string `json:"-"`
	Password          string `json:"password,omitempty"`
}

// Validate - валидация данных
func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Login, validation.Required, validation.Length(6, 100)), // проверка логина
		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(8, 100)), // проверка пароля
	)
}

// BeforeCreate - шифрование пароля и запись его в EncryptedPassword
func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = enc
	}
	return nil
}

// Sanitaze - очистка поля Password
func (u *User) Sanitaze() {
	u.Password = ""
}

// ComparePassword - сравнивание паролей, то шо послал юзер и то шо есть в бд
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

// encryptString - шифрование строки
func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
