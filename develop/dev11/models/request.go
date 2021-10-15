package models

import (
	"errors"
	"time"
)

var (
	errInvalidUserID = errors.New("значение user_id должно быть целым и положительным")
	errInvalidDate   = errors.New("дата должна быть в формате YYYY-MM-DD")
	errInvalidInfo   = errors.New("поле info обязательное и должно быть длиной как минимум в 3 символа")
)

// EventRequest играет роль промежуточного хранилища ещё не проверенных на корректность данных.
type EventRequest struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Date   string `json:"date"`
	Info   string `json:"info"`
}

// Структура EventRequest используется только для хранения значений параметров API-методов /create_event и /update_event, => метод Validate
// осуществляет валидацию значений параметров только этих двух API-методов. Валидация - проверка значений UserID, Date и Info на корректность
func (e *EventRequest) Validate() error {
	if e.UserID <= 0 {
		return errInvalidUserID
	}

	if _, err := time.Parse("2006-05-02", e.Date); err != nil {
		return errInvalidDate
	}
	if len(e.Info) == 0 {
		return errInvalidInfo
	}
	return nil
}

// NewEventFromRequest создаёт из валидированного значения EventRequest "окончательный" ивент, со стопроцентно корректными значениями полей
func NewEventFromRequest(e *EventRequest) *Event {
	return &Event{
		ID:     e.ID,
		UserID: e.UserID,
		Date:   e.Date,
		Info:   e.Info,
	}
}
