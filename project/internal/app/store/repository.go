package store

import "github.com/ENSLERMAN/soft-eng/project/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	FindByID(int) (*model.User, error)
	FindByLogin(string) (*model.User, error)
}

type BillRepository interface {
	CreateBill(*model.Bill, int) error
}
