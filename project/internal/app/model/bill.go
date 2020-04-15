package model

import (
	"math/rand"
	"time"
)

type Bill struct {
	ID      int `json:"id"`
	Type    int `json:"type"`
	Number  int `json:"number"`
	Balance int `json:"money"`
}

func RandomizeCardNumber() int {
	rand.Seed(time.Now().UnixNano())
	min := 1000000000000000
	max := 9999999999999999
	number := rand.Intn(max-min+1) + min
	return number
}
