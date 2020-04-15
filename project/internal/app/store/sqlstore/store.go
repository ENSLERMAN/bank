package sqlstore

import (
	"github.com/ENSLERMAN/soft-eng/project/internal/app/store"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Store struct {
	db             *sqlx.DB
	UserRepository *UserRepository
	BillRepository *BillRepository
}

func New(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.UserRepository != nil {
		return s.UserRepository
	}

	s.UserRepository = &UserRepository{store: s}

	return s.UserRepository
}

func (s *Store) Bill() store.BillRepository {
	if s.BillRepository != nil {
		return s.BillRepository
	}

	s.BillRepository = &BillRepository{store: s}

	return s.BillRepository
}
