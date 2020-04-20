package sqlstore

import (
	"github.com/ENSLERMAN/soft-eng/project/internal/app/model"
	"github.com/sirupsen/logrus"
)

type BillRepository struct {
	store *Store
}

func (r *BillRepository) CreateBill(u *model.Bill, id int) error {
	_, err := r.store.db.Exec(`INSERT INTO bank.bills 
		(type_bill, number_bill, balance)
		VALUES ($1, $2, $3) returning id`,
		&u.Type, &u.Number, 0)
	if err != nil {
		return err
	}

	_, err = r.store.db.Exec(`INSERT INTO bank.clients_bills 
		(bill_id, userid) VALUES ($1, $2)`, u.ID, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *BillRepository) DeleteBill(id int) error {
	_, err := r.store.db.Exec("DELETE FROM bank.clients_bills WHERE bill_id = $1", id)
	if err != nil {
		return err
	}
	_, err = r.store.db.Exec("DELETE FROM bank.bills WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *BillRepository) GetAllUserBills(id int) ([]*model.Bill, error) {

	arr := make([]*model.Bill, 0)

	rows, err := r.store.db.Queryx(`
		SELECT bills.id, number_bill, balance::numeric::float8, name FROM bank.clients_bills
    		INNER JOIN bank.bills
        		ON clients_bills.bill_id = bills.id
    		INNER JOIN bank.type_bill
        		ON bills.type_bill = type_bill.id
		WHERE clients_bills.userid = $1`, id)
	if err != nil {
		logrus.Error(err)
	}
	for rows.Next() {
		u := new(model.Bill)
		err := rows.Scan(
			&u.ID,
			&u.Number,
			&u.Balance,
			&u.Name,
		)
		if err != nil {
			return nil, err
		}
		arr = append(arr, u)
	}
	err = rows.Err()
	if err != nil {
		logrus.Error(err)
	}
	return arr, nil
}
