package model

import (
	"math/rand"
	"time"
)

type Bill struct {
	ID      int     `json:"id"`
	Type    int     `json:"type"`
	Name    string  `json:"name"`
	Number  int     `json:"number"`
	Balance float32 `json:"money"`
}

type ClientBill struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	BillID int `json:"bill_id"`
}

func RandomizeCardNumber() int {
	rand.Seed(time.Now().UnixNano())
	min := 1000000000000000
	max := 9999999999999999
	number := rand.Intn(max-min+1) + min
	return number
}
