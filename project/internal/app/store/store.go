package store

type Store interface {
	User() UserRepository
	Bill() BillRepository
}
