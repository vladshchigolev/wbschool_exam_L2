package cut

import (
	"dev06/managers"
	"errors"
	"log"
	"strings"
)

var (
	errDataNotProvided      = errors.New("Вы ничего не написали (ни одного столбца)")
	errInvalidFieldPosition = errors.New("Какие-то (или все) номера полей, указанные для вывода, не существуют")
)

// ManCut ...
type ManCut struct {
	manager *managers.ConsoleManager
	options Options  // Поле options содержит меняющиеся параметры (те, что могут устанавливаться с помощью флагов), используем композицию для добавления необходимых полей
	data    []string // входные данные
	result  []string // результат работы программы
	//delimeter              string // значение разделителя
	//onlySeparated          bool // только строки с разделителем
}

// New ...
func New(manager *managers.ConsoleManager) *ManCut { // Конструктор, создаёт новое значение ManCut и возвращает указатель на него,
	return &ManCut{ // при этом на вход New принимает значение типа, соответствующего типу ConsoleMananger
		manager: manager, // Инициализируется ненулевым значением только поле manager
	}
}

// Метод объекта ManCut (*ManCut) - ApplyOptions
// Здесь сначала инициализируем структуру Options дефолтными значениями, затем передаём его функциям типа Option, изменяющим один из 3-х параметров КАЖДАЯ
func (m *ManCut) ApplyOptions(options ...Option) *ManCut { // ApplyOptions принимает слайс функций с сигнатурой "func(*Options) error" ...
	opts := GetDefaultOptions()   // Возвращает структуру Options, инициализированную по дефолту, присваиваем это значение opts
	for _, opt := range options { // Каждой ф-ции, что лежит в слайсе нужно дать на вход указатель на структуру Option
		if opt != nil { // Проверяем значение opt на каждой итерации на равенство nil, если не nil, вызываем функцию с параметром типа *Options (указателем на структуру)
			if err := opt(&opts); err != nil { // ...а затем вызывает эти функции в цикле, передавая им указатель на структуру Options (заполненную дефолтными значениями), созданную в строке 36
				log.Fatal(err) // Если вызванная функция вернула не nil, завершаемся с ошибкой
			}
		}
	}
	m.options = opts // Присваиваем полю options экземпляра ManCut измененную (или нет) структуру opts
	return m
}

// Cut разбивает по разделителю
func (m *ManCut) Cut() error {
	data, err := m.manager.Read() // Метод Read значения ConsoleManager читает пользовательский ввод построчно, добавляет каждую строку в слайс
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errDataNotProvided // Если метод Read вернул нам пустой слайс строк, вернём ошибку
	}
	m.data = data // Далее присваиваем прочитанные данные полю data структуры ManCut

	SelectFields(m)
	return nil
}

func (m *ManCut) OutputResult() error {
	if err := m.manager.Write(m.result); err != nil {
		return err
	}
	return nil
}
func SelectFields(m *ManCut) {
	for _, line := range m.data { // line здесь - это строка пользовательского ввода
		if !strings.Contains(line, m.options.delimeter) { // Если в этой строке нет разделителя...
			if !m.options.separated { // ...и не указано, что выводим строки только с разделителем...
				m.result = append(m.result, line) // ...то добавляем строку к результату
			}
			continue
		}
		samples := strings.Split(line, m.options.delimeter) // Если же в строке разделитель есть, то делим её по нему на подстроки (это будут столбцы)

		var prepared []string
		for _, fieldID := range m.options.fields { // Перебираем номера столбцов, которые нужно вывести...
			if fieldID > len(samples)-1 { // Если какие-то из номеров полей, указанные пользователем, не существуют, завершаемся с ошибкой
				log.Fatal(errInvalidFieldPosition)
			}
			prepared = append(prepared, samples[fieldID]) // ...и заполняем новый слайс элементами с нужным индексом из слайса samples
		}
		preparedString := strings.Join(prepared, m.options.delimeter)
		m.result = append(m.result, preparedString)
	}
}
