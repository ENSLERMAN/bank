package store

import "github.com/ENSLERMAN/soft-eng/project/internal/app/model"

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) (*model.User, error) {
	if err := u.Validate(); err != nil {
		return nil, err
	}

	if err := u.BeforeCreate(); err != nil {
		return nil, err
	}

	if err := r.store.db.QueryRowx(`INSERT INTO bank.clients 
		(login, password, name, surname, patronymic, passport)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		u.Login, u.EncryptedPassword, u.Name, u.Surname, u.Patronymic, u.Passport,
	).Scan(&u.ID); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) FindByLogin(login string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRowx(
		`SELECT id, login, surname, name, patronymic FROM bank.clients WHERE login = $1`,
		login,
	).Scan(
		&u.ID,
		&u.Login,
		&u.Surname,
		&u.Name,
		&u.Patronymic,
	); err != nil {
		return nil, err
	}
	return u, nil
}
