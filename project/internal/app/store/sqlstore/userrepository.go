package sqlstore

import (
	"database/sql"
	"github.com/ENSLERMAN/soft-eng/project/internal/app/model"
	"github.com/ENSLERMAN/soft-eng/project/internal/app/store"
)

type UserRepository struct {
	store *Store
}

// Create - метод для создания юзера
func (r *UserRepository) Create(u *model.User) error {
	// проверяем валидность данных
	if err := u.Validate(); err != nil {
		return err
	}

	// шифруем данные
	if err := u.BeforeCreate(); err != nil {
		return err
	}

	// создаем клиента и номер счета по умолчанию
	b := &model.Bill{}
	number := model.RandomizeCardNumber()

	if err := r.store.db.QueryRowx(`INSERT INTO bank.clients 
		(login, password, name, surname, patronymic, passport)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		u.Login, u.EncryptedPassword, u.Name, u.Surname, u.Patronymic, u.Passport,
	).Scan(&u.ID); err != nil {
		return err
	}

	/* типы счетов
	1 - Default Bill
	2 - MasterCard
	3 - Visa
	4 - Mir
	*/
	// по дефолту создается: "Default bill" с нулем на балансе
	if err := r.store.db.QueryRowx(`INSERT INTO bank.bills 
		(type_bill, number, balance)
		VALUES ($1, $2, $3) RETURNING id`,
		1, number, 0,
	).Scan(&b.ID); err != nil {
		return err
	}


	_, err := r.store.db.Exec(`INSERT INTO bank.clients_bills 
		(bill_id, user_id)
		VALUES ($1, $2)`,
		&b.ID, &u.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// FindByLogin - метод поиска юзера по его логину
func (r *UserRepository) FindByLogin(login string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRowx(
		`SELECT id, login, surname, name, patronymic, passport, password FROM bank.clients WHERE login = $1`,
		login,
	).Scan(
		&u.ID,
		&u.Login,
		&u.Surname,
		&u.Name,
		&u.Patronymic,
		&u.Passport,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

// FindByID - метод поиска юзера по его ID
func (r *UserRepository) FindByID(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRowx(
		`SELECT id, login, surname, name, patronymic, passport, password FROM bank.clients WHERE id = $1`,
		id,
	).Scan(
		&u.ID,
		&u.Login,
		&u.Surname,
		&u.Name,
		&u.Patronymic,
		&u.Passport,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}
