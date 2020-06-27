package sqlstore

import (
	"github.com/ENSLERMAN/soft-eng/project/internal/app/model"
	"github.com/ENSLERMAN/soft-eng/project/internal/app/store"
	"github.com/sirupsen/logrus"
	"log"
)

type BillRepository struct {
	// отдаем ссылку на store.
	store *Store
}

// CreateBill - метод создания счета.
func (r *BillRepository) CreateBill(u *model.Bill, id int) error {
	// кидаем данные в бд
	if err := r.store.db.QueryRowx(`INSERT INTO bank.bills 
		(type_bill, number, balance, user_id)
		VALUES ($1, $2, $3, $4) returning id`, &u.Type, &u.Number, 0, id).Scan(&u.ID); err != nil {
		return err
	}

	return nil
}

// DeleteBill - метод закрытия счета ( удаления ).
func (r *BillRepository) DeleteBill(id int) error {
	_, err := r.store.db.Exec("DELETE FROM bank.bills WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

// GetAllUserBills - метод для просмотра всех счетов юзера.
func (r *BillRepository) GetAllUserBills(id int) ([]*model.Bill, error) {

	arr := make([]*model.Bill, 0)

	rows, err := r.store.db.Queryx(`
		SELECT id, number, balance::numeric::float8, type_bill, user_id FROM bank.bills
		WHERE bills.user_id = $1`, id)
	if err != nil {
		logrus.Error(err)
	}
	for rows.Next() {
		u := new(model.Bill)
		err := rows.Scan(
			&u.ID,
			&u.Number,
			&u.Balance,
			&u.Type,
			&u.UserID,
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

// GetAllUserPayments получить историю переводов пользователя
func (r *BillRepository) GetAllUserPayments(id int) ([]*model.Payment, error) {

	arr := make([]*model.Payment, 0)

	rows, err := r.store.db.Queryx(`
		SELECT DISTINCT id, sender, recipient, amount::numeric::float8, 
		to_char(time AT TIME ZONE 'Europe/Moscow', 'HH24:MI:SS'),
		to_char(time AT TIME ZONE 'Europe/Moscow', 'DD.MM.YYYY'),
		sender_id, rec_id
		FROM bank.payments
		WHERE (sender_id = $1) OR (rec_id = $1)`, id)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()

	for rows.Next() {
		u := new(model.Payment)
		err := rows.Scan(
			&u.ID,
			&u.Sender,
			&u.Recipient,
			&u.Amount,
			&u.Time,
			&u.Date,
			&u.SenderID,
			&u.RecID,
		)
		if err != nil {
			return nil, err
		}
		//fmt.Printf("sender: %d = %t , rec: %d = %t\n", u.Sender, r.FindByBill(id, u.Sender), u.Recipient, r.FindByBill(id, u.Recipient))
		if u.SenderID == id && u.RecID != id {
			u.Type = 1
		} else if u.SenderID != id && u.RecID == id {
			u.Type = 2
		} else if u.SenderID == id && u.RecID == id {
			u.Type = 3
		}

		arr = append(arr, u)
	}

	return arr, nil
}

// FindByUser - метод для сопоставления номера счета и юзера, нужен для перевода денег.
func (r *BillRepository) FindByUser(userID, billID int) bool {

	var count int

	err := r.store.db.QueryRowx(`
		SELECT COUNT(id) FROM bank.bills WHERE id = $1 and user_id = $2`, billID, userID,
	).Scan(
		&count,
	); if err != nil {
		return false
	}

	if count == 1 {
		return true
	}

	return false
}

// GetUserBillByID - получить счет клиента по ид
func (r *BillRepository) GetUserBillByID(id int) (*model.Bill, error) {

	u := &model.Bill{}

	err := r.store.db.QueryRowx(`
		SELECT id, type_bill, number, balance::numeric::float8, user_id FROM bank.bills WHERE bills.id = $1`, id,
	).Scan(
		&u.ID,
		&u.Type,
		&u.Number,
		&u.Balance,
		&u.UserID,
	); if err != nil {
		return nil, err
	}
	if err != nil {
		logrus.Error(err)
	}
	return u, nil

}

func (r *BillRepository) FindByBill(userID, number int) bool {

	var nuSho bool

	rows, err := r.store.db.Queryx(`
		SELECT exists(select * FROM bank.bills WHERE (number = $1 and user_id = $2))`, number, userID)
	if err != nil {
		return false
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&nuSho)
		if err != nil {
			log.Fatal(err)
			return false
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return false
	}

	return nuSho
}

// TransferMoney - метод перевода денег.
func (r *BillRepository) TransferMoney(NumberDest int, Amount uint, billID int) error {

	u := &model.Bill{}
	k := &model.Bill{}
	var reqID int
	var resID int

	req := &u.Balance // баланс отправителя.
	res := &k.Balance // баланс получателя.
	var NumberSender int // номер карты отправителя.

	// получаем номер карты и баланс отправителя.
	if err := r.store.db.QueryRowx(`
		SELECT user_id, balance::numeric::float8, number from bank.bills WHERE id = $1`, billID,
	).Scan(
		&reqID,
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
		SELECT user_id, balance::numeric::float8 from bank.bills WHERE number = $1`, NumberDest,
	).Scan(
		&resID,
		&res,
	); err != nil {
		return err
	}

	// обновляем данные у получателя и отправителя.
	tx, err := r.store.db.Begin()

	_, err = tx.Exec(`UPDATE bank.bills SET balance = $1 WHERE id = $2`, *req-float32(Amount), billID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = r.store.db.Exec(`UPDATE bank.bills SET balance = $1 WHERE number = $2`, *res+float32(Amount), NumberDest)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = r.store.db.Exec(`INSERT INTO bank.payments
		(sender, recipient, amount, time, sender_id, rec_id) 
		VALUES ($1, $2, $3, 'now', $4, $5)`, NumberSender, NumberDest, Amount, reqID, resID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if tx.Commit() != nil {
		return err
	}

	return nil
}

func (r *BillRepository) GetUserPaymentsByID(userID, id int) ([]*model.Payment, error) {

	arr := make([]*model.Payment, 0)
	var number int

	err := r.store.db.QueryRowx(`SELECT number from bank.bills WHERE id = $1`, id).Scan(
		&number,
	); if err != nil {
		return nil, err
	}

	rows, err := r.store.db.Queryx(`
		SELECT DISTINCT id, sender, recipient, amount::numeric::float8, 
		to_char(time AT TIME ZONE 'Europe/Moscow', 'HH24:MI:SS'),
		to_char(time AT TIME ZONE 'Europe/Moscow', 'DD.MM.YYYY'),
		sender_id, rec_id
		FROM bank.payments
		WHERE (sender = $1 OR recipient = $1)`, number)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()

	for rows.Next() {
		u := new(model.Payment)
		err := rows.Scan(
			&u.ID,
			&u.Sender,
			&u.Recipient,
			&u.Amount,
			&u.Time,
			&u.Date,
			&u.SenderID,
			&u.RecID,
		)
		if err != nil {
			return nil, err
		}
		// t f
		if u.SenderID == userID && u.RecID != userID {
			u.Type = 1
		// f t
		} else if u.SenderID != userID && u.RecID == userID {
			u.Type = 2
		// t t
		} else if u.RecID == u.SenderID {
			u.Type = 3
		}

		arr = append(arr, u)
	}

	return arr, nil
}

func (r *BillRepository) GetMoney(id int) error {

	var NumberDest int
	var UserID int
	var balance []uint8

	// обновляем данные у получателя и отправителя.
	tx, err := r.store.db.Begin()

	// получаем номер карты и баланс отправителя.
	if err := r.store.db.QueryRowx(`
		SELECT number, balance, user_id from bank.bills WHERE id = $1`, id,
	).Scan(
		&NumberDest,
		&balance,
		&UserID,
	); err != nil {
		return err
	}

	_, err = r.store.db.Exec(`UPDATE bank.bills SET balance = balance + money(5000) WHERE id = $1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = r.store.db.Exec(`INSERT INTO bank.payments
		(sender, recipient, amount, time, sender_id, rec_id) 
		VALUES ($1, $2, $3, 'now', $4, $5)`, 1000000000001001, NumberDest, 5000, 0, UserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if tx.Commit() != nil {
		return err
	}

	return nil

}


func (r *BillRepository) GetRestOfTheBills(id int) ([]int, error) {

	arr := make([]int, 0)
	var Number int

	rows, err := r.store.db.Queryx(`
		SELECT number FROM bank.bills WHERE bills.user_id != $1`, id)
	if err != nil {
		logrus.Error(err)
	}
	for rows.Next() {
		err := rows.Scan(
			&Number,
		)
		if err != nil {
			return nil, err
		}
		arr = append(arr, Number)
	}
	err = rows.Err()
	if err != nil {
		logrus.Error(err)
	}
	return arr, nil

}
