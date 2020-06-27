package teststore

import "github.com/ENSLERMAN/soft-eng/project/internal/app/model"

type BillRepository struct {
	store *Store
	bills map[int]*model.Bill
}

func (b BillRepository) CreateBill(*model.Bill, int) error {
	return nil
}

func (b BillRepository) GetAllUserBills(int) ([]*model.Bill, error) {
	return nil, nil
}

func (b BillRepository) FindByUser(int, int) (*model.ClientBill, error) {
	return nil, nil
}

func (b BillRepository) TransferMoney(int, uint, int) error {
	return nil
}

func (b BillRepository) DeleteBill(int) error {
	return nil
}

