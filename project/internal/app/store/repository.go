package store

import "github.com/ENSLERMAN/soft-eng/project/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	FindByID(int) (*model.User, error)
	FindByLogin(string) (*model.User, error)
}

type BillRepository interface {
	CreateBill(*model.Bill, int) error
	GetAllUserBills(int) ([]*model.Bill, error)
	FindByUser(int, int) (*model.ClientBill, error)
	TransferMoney(int, uint, int) error
	DeleteBill(int) error
}
