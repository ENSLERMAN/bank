package sqlstore

import (
	"database/sql"
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

// todo: ПЕРЕПИШИ ЭТУ ХУЙНЮ, ХОСПАДЕ, ТЫ ШО ДАУН?
func (r *BillRepository) GetAllUserPayments(id int) ([]*model.Payment, error) {

	arr := make([]*model.Payment, 0)

	rows, err := r.store.db.Queryx(`
		SELECT DISTINCT id, sender, recipient, amount::numeric::float8, 
		to_char(time AT TIME ZONE 'Europe/Moscow', 'HH24:MI:SS'),
		to_char(time AT TIME ZONE 'Europe/Moscow', 'DD.MM.YYYY')
		FROM bank.payments`)
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
		)
		if err != nil {
			return nil, err
		}
		//fmt.Printf("sender: %d = %t , rec: %d = %t\n", u.Sender, r.FindByBill(id, u.Sender), u.Recipient, r.FindByBill(id, u.Recipient))
		if r.FindByBill(id, u.Sender) == false && r.FindByBill(id, u.Recipient) == false {
			continue
		} else if r.FindByBill(id, u.Sender) == true && r.FindByBill(id, u.Recipient) == false {
			u.Type = 1
		} else if r.FindByBill(id, u.Sender) == false && r.FindByBill(id, u.Recipient) == true {
			u.Type = 2
		} else if r.FindByBill(id, u.Sender) == true && r.FindByBill(id, u.Recipient) == true {
			u.Type = 3
		}

		arr = append(arr, u)
	}

	return arr, nil
}

// FindByUser - метод для сопоставления номера счета и юзера, нужен для перевода денег.
func (r *BillRepository) FindByUser(userID, billID int) bool {

	_, err := r.store.db.Exec(`
		SELECT id, user_id FROM bank.bills WHERE id = $1 and user_id = $2`, billID, userID)
	if err != nil || err == sql.ErrNoRows {
		return false
	}

	return true
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
		(sender, recipient, amount, time) VALUES ($1, $2, $3, 'now')`, NumberSender, NumberDest, Amount)
	if err != nil {
		tx.Rollback()
		return err
	}

	if tx.Commit() != nil {
		return err
	}

	return nil
}
