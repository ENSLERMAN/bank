// Пакет store - это хранилище для моделей и их методов.
// sqlstore - это хранилище для работы с реальной бд.
// teststore - это хранилище для тестов с фейк бд.
package store

// Store - хранилище для моделей и их методов.
type Store interface {
	User() UserRepository
	Bill() BillRepository
}
