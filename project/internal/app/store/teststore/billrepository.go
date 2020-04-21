package teststore

import "github.com/ENSLERMAN/soft-eng/project/internal/app/model"

type BillRepository struct {
	store *Store
	bills map[int]*model.Bill
}
