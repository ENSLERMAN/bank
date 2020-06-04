package store

import "github.com/ENSLERMAN/soft-eng/project/internal/app/model"

// UserRepository содержит методы для модели User.
type UserRepository interface {
	Create(*model.User) error
	FindByID(int) (*model.User, error)
	FindByLogin(string) (*model.User, error)
}

// BillRepository содержит методы для модели Bill.
type BillRepository interface {
	CreateBill(*model.Bill, int) error
	GetAllUserBills(int) ([]*model.Bill, error)
	GetAllUserPayments(int) ([]*model.Payment, error)
	FindByUser(int, int) bool
	FindByBill(int, int) bool
	TransferMoney(int, uint, int) error
	DeleteBill(int) error
}
