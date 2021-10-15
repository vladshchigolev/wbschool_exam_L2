package clock

import (
	"fmt"
	"github.com/beevik/ntp"
	"time"
)

const (
	// Константа для обращения к стандартному серверу времени
	DefaultHost = "0.beevik-ntp.pool.ntp.org"
	AlternativeHost = "www.example.com"
)


type IClock interface { // Тип интерфейса, в котором определены методы
	CurrentTime() (time.Time, time.Time)
	SetHost(string) error
	String() string
}

// Clock - базовые часы
type Clock struct {
	response *ntp.Response // Является ли это композицией? Поле response имеет тип указателя на структуру Response
	host     string        // Поля структуры неэкспортируемые (доступны ТОЛЬКО в пределах этого пакета) - таким образом
						   // осуществляем сокрытие реализации (инкапсуляцию). Для создания экземпляра типа за пределами пакета
}						   // Clock используем "API" - конструктор (функция New)

// Конструктор, возвращающий ЗНАЧЕНИЕ ИНТЕРФЕЙСА (в нашем случае динамическим типом интерфейса будет тип *Clock, т. к. именно он
// поддерживает/реализует этот интерфейс - то есть имеет все методы, необходимые для его поддержки) и ЗНАЧЕНИЕ ОШИБКИ
func New(host string) (IClock, error) { // вход: имя хоста, выход: функция вернет только то, что соответствует интерфейсу IClock
	response, err := ntp.Query(host) // Query возвращает указатель на значение типа Response (структура), содержащее текущее время и различные дополнительные метаданные
	if err != nil {					 // Также Query возвращает значение ошибки, обрабатываем её
		return nil, err // Если значение ошибки не nil, то в качестве результата возвращаем nil, а в качестве ошибки значение ошибки (не nil)
	}
	return &Clock{ // Возвращаем указатель на Clock (именно тип *Clock является получателем всех необходимых для реализации интерфейса IClock методов)
		response: response, // В структуру Clock добавляем значение-указатель на возвращенную Query структуру
		host:     host,
	}, nil
}

// CurrentTime возвращает ТОЧНОЕ и ЛОКАЛЬНОЕ время
func (c *Clock) CurrentTime() (time.Time, time.Time) {
	prec := time.Now().Add(c.response.ClockOffset)
	loc := time.Now()
	return prec, loc
}

// SetHost ...
func (c *Clock) SetHost(host string) error { // Когда экземпляр создан, а НАДО изменить адрес хоста,
	response, err := ntp.Query(host)         // метод SetHost позволяет перезаписать поля нашего экземпляра
	if err != nil {							 // новыми значениями. Снова выполняем запрос Query() с новым именем хоста в кач-ве параметра,
		return err							 // перезаписываем поле response нашего экземпляра результатом этого запроса. Также перезаписываем
	}					 					 // поле с именем хоста
	c.response = response
	return nil
}

// String ...
func (c *Clock) String() string {
	prec, cur := c.CurrentTime()
	return fmt.Sprintf("Precise:%v\nLocal:%v", prec, cur) // %v - произвольное значение (подходящий формат выбирается на основании типа передаваемого значения)
}