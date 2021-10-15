package managers

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	errConsoleInput = errors.New("Не удалось прочитать данные из ст. ввода")
)

// ConsoleManager пишет/читает В и ИЗ stdout/stdin
// ConsoleManager (а точнее *ConsoleManager) имеет методы Read и Write
type ConsoleManager struct { // Сущность ConsoleManager имеет поля интерфейсного типа, значит могут принимать значения любых соответствующих
	writer io.Writer		 // этим интерфейсам типов (имеющих методы Write и Read соответственно)
	reader io.Reader
}

// NewConsoleManager  - конструктор, создает значение ConsoleManager, возвращает указатель на него
func NewConsoleManager(reader io.Reader, writer io.Writer) *ConsoleManager { // Когда будем вызывать эту функцию, передадим ей os.Stdin, os.Stdout,
	return &ConsoleManager{                                                  // имеющие методы Read/Write соответственно
		writer: writer,
		reader: reader,
	}
}

// Read читает из stdin, возвращает слайс строк, где каждый элемент - строка пользовательского ввода
func (cm *ConsoleManager) Read() ([]string, error) { // Получатель метода - указатель на значение ConsoleManager
	data := make([]string, 0)
	scanner := bufio.NewScanner(cm.reader) // Создать "буферизованное средство чтения"
	for {
		scanner.Scan() // Читаем построчно...
		text := scanner.Text()
		if len(text) != 0 {
			data = append(data, text) // Каждое прочитанное строковое значение добавляем в слайс
		} else { // ...до тех пор, пока очередной строкой не будет прочитана ПУСТАЯ строка
			break
		}
	}
	if scanner.Err() != nil {
		return nil, errConsoleInput
	}
	return data, nil // Этот метод может вернуть в кач-ве значения data пустую строку, => надо будет обработать ошибку при его вызове
}

// Write выводит результат (значения поля result структуры ManCut) в stdout
func (cm *ConsoleManager) Write(data []string) error {
	_, err := fmt.Fprintln(cm.writer, strings.Join(data, "\n")) // В нашем случае работаем с stdout (cm.writer это io.stdout)
	return err
}
