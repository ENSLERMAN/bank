package sqlstore

import "github.com/ENSLERMAN/soft-eng/project/internal/app/model"

type BillRepository struct {
	store *Store
}

func (r *BillRepository) CreateBill(u *model.Bill, id int) error {
	if err := r.store.db.QueryRowx(`INSERT INTO bank.bills 
		(type_bill, number_bill, balance)
		VALUES ($1, $2, $3) RETURNING id`,
		u.Type, u.Number, 0,
	).Scan(&u.ID); err != nil {
		return err
	}

	return r.store.db.QueryRowx(`INSERT INTO bank.clients_bills 
		(bill_id, userid)
		VALUES ($1, $2)`,
		&u.ID, id,
	).Scan()
}
