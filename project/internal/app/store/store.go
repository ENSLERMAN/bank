package store

// Store - хранилище для моделей и их методов
type Store interface {
	User() UserRepository
	Bill() BillRepository
}
