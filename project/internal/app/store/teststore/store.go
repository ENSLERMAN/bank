package teststore

import (
	"github.com/ENSLERMAN/soft-eng/project/internal/app/model"
	"github.com/ENSLERMAN/soft-eng/project/internal/app/store"
)

type Store struct {
	UserRepository *UserRepository
	BillRepository *BillRepository
}

func (s *Store) Bill() store.BillRepository {
	if s.BillRepository != nil {
		return s.BillRepository
	}

	s.BillRepository = &BillRepository{
		store: s,
		bills: make(map[int]*model.Bill),
	}

	return s.BillRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepository {
	if s.UserRepository != nil {
		return s.UserRepository
	}

	s.UserRepository = &UserRepository{
		store: s,
		users: make(map[int]*model.User),
	}

	return s.UserRepository
}
