package store

import "dev11/models"

// Структура Store хранит все ивенты.
// Значение Store имеет ссылку на значение EventRepository. EventRepository, в свою очередь, имеет ссылку на значение Store.
// База данных не должна реализовывать функциональность вроде Create/Update/Delete Event и т. д. Делегируем её отдельному
// классу (типу) EventRepository
type Store struct {
	db         map[int]*models.Event // Значения ключа - id-шники ивентов (уникальны для каждого)
	repository *EventRepository
}

// New ...
func New() *Store {
	return &Store{}
}

// Open создаёт мапу, присваивает это значение полю db
func (s *Store) Open() error {
	db := make(map[int]*models.Event)
	s.db = db
	return nil
}

// Close заглушка
func (s *Store) Close() error {
	return nil
}

// EventRepository ...
func (s *Store) EventRepository() *EventRepository {
	if s.repository == nil {
		s.repository = &EventRepository{store: s}
	}
	return s.repository
}
