package sqlstore

import (
	"database/sql"
	"github.com/ENSLERMAN/soft-eng/project/internal/app/model"
	"github.com/ENSLERMAN/soft-eng/project/internal/app/store"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	b := &model.Bill{}
	number := model.RandomizeCardNumber()

	if err := r.store.db.QueryRowx(`INSERT INTO bank.clients 
		(login, password, name, surname, patronymic, passport)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		u.Login, u.EncryptedPassword, u.Name, u.Surname, u.Patronymic, u.Passport,
	).Scan(&u.ID); err != nil {
		return err
	}

	if err := r.store.db.QueryRowx(`INSERT INTO bank.bills 
		(type_bill, number_bill, balance)
		VALUES ($1, $2, $3) RETURNING id`,
		1, number, 0,
	).Scan(&b.ID); err != nil {
		return err
	}

	return r.store.db.QueryRowx(`INSERT INTO bank.clients_bills 
		(bill_id, userid)
		VALUES ($1, $2)`,
		&b.ID, &u.ID,
	).Scan()
}

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
