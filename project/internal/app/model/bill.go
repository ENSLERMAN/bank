// Пакет model содержит все модели данных
package model

import (
	"math/rand"
	"time"
)

// Bill - структура счета.
type Bill struct {
	ID      int     `json:"id"`
	Type    int     `json:"type"`
	Number  int     `json:"number"`
	Balance float32 `json:"money"`
	UserID int `json:"user_id"`
}

type Payment struct {
	ID int `json:"id"`
	Sender int `json:"sender"`
	Recipient int `json:"recipient"`
	Amount float32 `json:"amount"`
	Time string `json:"time"`
	Date string `json:"date"`
}

// RandomizeCardNumber - рандомим номер карты для создания счета.
func RandomizeCardNumber() int {
	rand.Seed(time.Now().UnixNano())
	min := 1000000000000000
	max := 9999999999999999
	number := rand.Intn(max-min+1) + min
	return number
}
