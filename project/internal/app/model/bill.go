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
	Name    string  `json:"name"`
	Number  int     `json:"number"`
	Balance float32 `json:"money"`
}

// ClientBill - структура таблицы в бд client_bills.
type ClientBill struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
	BillID int `json:"bill_id"`
}

// RandomizeCardNumber - рандомим номер карты для создания счета.
func RandomizeCardNumber() int {
	rand.Seed(time.Now().UnixNano())
	min := 1000000000000000
	max := 9999999999999999
	number := rand.Intn(max-min+1) + min
	return number
}
