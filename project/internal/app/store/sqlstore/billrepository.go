package sqlstore

import (
	"database/sql"
	"github.com/ENSLERMAN/soft-eng/project/internal/app/model"
	"github.com/ENSLERMAN/soft-eng/project/internal/app/store"
	"github.com/sirupsen/logrus"
)

type BillRepository struct {
	// отдаем ссылку на store.
	store *Store
}

// CreateBill - метод создания счета.
func (r *BillRepository) CreateBill(u *model.Bill, id int) error {
	// кидаем данные в бд
	if err := r.store.db.QueryRowx(`INSERT INTO bank.bills 
		(type_bill, number, balance)
		VALUES ($1, $2, $3) returning id`, &u.Type, &u.Number, 0).Scan(&u.ID); err != nil {
		return err
	}

	_, err := r.store.db.Exec(`INSERT INTO bank.clients_bills 
		(bill_id, user_id) VALUES ($1, $2)`, &u.ID, id)
	if err != nil {
		return err
	}

	return nil
}

// DeleteBill - метод закрытия счета ( удаления ).
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

// GetAllUserBills - метод для просмотра всех счетов юзера.
func (r *BillRepository) GetAllUserBills(id int) ([]*model.Bill, error) {

	arr := make([]*model.Bill, 0)

	rows, err := r.store.db.Queryx(`
		SELECT bills.id, number, balance::numeric::float8, name FROM bank.clients_bills
    		INNER JOIN bank.bills
        		ON clients_bills.bill_id = bills.id
    		INNER JOIN bank.type_bill
        		ON bills.type_bill = type_bill.id
		WHERE clients_bills.user_id = $1`, id)
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

// FindByUser - метод для сопоставления номера счета и юзера, нужен для перевода денег.
func (r *BillRepository) FindByUser(userID, billID int) (*model.ClientBill, error) {
	u := &model.ClientBill{}
	if err := r.store.db.QueryRowx(
		`SELECT id, user_id, bill_id FROM bank.clients_bills WHERE user_id = $1 and bill_id = $2`,
		userID, billID,
	).Scan(
		&u.ID,
		&u.UserID,
		&u.BillID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

// TransferMoney - метод перевода денег.
func (r *BillRepository) TransferMoney(NumberDest int, Amount uint, billID int) error {

	u := &model.Bill{}
	k := &model.Bill{}

	req := &u.Balance // баланс отправителя.
	res := &k.Balance // баланс получателя.
	var NumberSender int // номер карты отправителя.

	// получаем номер карты и баланс отправителя.
	if err := r.store.db.QueryRowx(`
		SELECT balance::numeric::float8, number from bank.bills WHERE id = $1`, billID,
	).Scan(
		&req,
		&NumberSender,
	); err != nil {
		return err
	}

	// если номер карты отправителя и получателя совпадают ретёрним ошибку.
	if NumberSender == NumberDest {
		return store.NumberSenderAndDest
	}

	// если сумма перевода больше того, что лежит на карте, возвращаем ошибку.
	if *req < float32(Amount) {
		return store.GreaterAmount
	}

	// получаем баланс получателя.
	if err := r.store.db.QueryRowx(`
		SELECT balance::numeric::float8 from bank.bills WHERE number = $1`, NumberDest,
	).Scan(
		&res,
	); err != nil {
		return err
	}

	// обновляем данные у получателя и отправителя.
	_, err := r.store.db.Exec(`UPDATE bank.bills SET balance = $1 WHERE id = $2`,
		*req-float32(Amount), billID)
	if err != nil {
		return err
	}

	_, err = r.store.db.Exec(`UPDATE bank.bills SET balance = $1 WHERE number = $2`,
		*res+float32(Amount), NumberDest)
	if err != nil {
		return err
	}

	return nil
}
