package store

import (
	"dev11/models"
	"errors"
)

// EventRepository ...
type EventRepository struct {
	store *Store
}

var (
	BaseTimeSample        = "2006-05-02"
	errEventAlreadyExists = errors.New("запись с таким id уже существует")
	errEventDoesNotExists = errors.New("такая запись не существует")
)

// Метод CreateEvent сохраняет переданный ему ивент в базу
func (e *EventRepository) CreateEvent(event *models.Event) error {
	if !e.checkIfExists(event) { // Если ивент с таким id-шником ещё не существует в базе,
		id := len(e.store.db) + 1 // генерируем для нового ивента свой id-шник,
		event.ID = id // присваиваем его полю ID
		e.store.db[id] = event // записываем ивент в мапу по ключу - id-шнику
		return nil
	}
	return errEventAlreadyExists // Если выражение !e.checkIfExists(event) == false, значит ивент с таким айдишником уже есть в базе
}

// UpdateEvent обновляет ивент в базе (мапе) по id-шнику (ключу)
func (e *EventRepository) UpdateEvent(event *models.Event) error {
	if e.checkIfExists(event) { // Если ивент с таким id-шником существует в мапе,
		//for id := range e.store.db {
		//	if id == event.ID {
		e.store.db[event.ID] = event
		return nil
		//}
		//}
	}
	return errEventDoesNotExists
}

// DeleteEvent удаляет ивент из базы (мапы) по id-шнику (ключу)
func (e *EventRepository) DeleteEvent(id int) error {
	_, ok := e.store.db[id]
	if !ok {
		return errEventDoesNotExists
	}
	delete(e.store.db, id)
	return nil
}

// GetEventsForDates получает ивент/ивенты из базы, попадающие в определенный временной диапазон и возвращает слайс указателей с ним/ними
// Мы будем использовать этот метод для диапазонов "день", "неделя", "месяц"
func (e *EventRepository) GetEventsForDates(startDate, endDate string) ([]*models.Event, error) {
	var events []*models.Event // Результирующий слайс ивентов
	for _, val := range e.store.db { // Перебираем все ивенты из базы
		if val.Date >= startDate && val.Date <= endDate {
			events = append(events, val) // Заполняем слайс ивентами, где значение поля Date попадает в заданный диапазон
		}
	}
	return events, nil
}
// Проверяет наличие ивента с таким id-шником в базе
func (e *EventRepository) checkIfExists(event *models.Event) bool {
	// Т. к. необходимо идентифицировать каждый отдельный ивент, пусть каждый ивент в базе имеет уникальный номер

	for id := range e.store.db { // Перебираем ключи мапы
		if id == event.ID {
			return true
		}
	}
	return false
}
