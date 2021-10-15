package cut

import (
	"errors"
	"strconv"
	"strings"
)

var (
	errUnparsedField      = errors.New("Не удалось распарсить указанные поля, используйте ТОЛЬКО целые числа, разделенные запятыми")
	errNegativeFieldValue = errors.New("Значение, представляющее номер столбца, должно быть НЕОТРИЦАТЕЛЬНЫМ")
	errEmptyFields        = errors.New("Вы не указали ни одного номера столбца для вывода")
)

type Options struct { // Структура для изменяющихся опций
	fields    []int // Выбранные для вывода колонки
	delimeter string // По какому разделителю разбиваем на колонки
	separated bool // Только строки с разделителем
}

func GetDefaultOptions() Options { // Функция возвращает нам экземпляр структуры Options, с полями, заполненными значениями по умолчанию
	return Options{
		fields:    []int{}, // Номера столбцов, которые необходимо вывести
		delimeter: "\t", // В кач-ве разделителя - TAB (по условию задания)
		separated: false,
	}
}

var DefaultOptions = GetDefaultOptions() // Присваиваем DefaultOptions экземпляр Options, возвращённый ф-ей GetDefaultOptions()

type Option func(*Options) error // Значением типа Option будет функция с соответствующей сигнатурой

// Далее идут 3 функции установки конкретной опции (т. е. пОля в структуре Options)
// Функция, возвращаемая SetFieldsOption, преобразует строку (пользовательский ввод, перечень номеров отбираемых столбцов) к слайсу значений int и присваивает его полю структуры Options
func SetFieldsOption(fields string) Option { // SetFieldsOption возвращает значение типа Option (функцию с сигнатурой func(*Options) error)
	return func(o *Options) error { // вход: указатель на значение Options ()
		resultedFields := []int{} // Этот слайс будем заполнять значениями, представляющими номера отбираемых для вывода столбцов
		if len(fields) == 0 { // Так как переданное ф-ции SetFieldsOption при вызове значение - строка, введённая пользователем (перечисление номеров столбцов для выбора), то если на вход SetFieldsOption дали пустую строку...
			return errEmptyFields // ...вернуть ошибку
		} // Если длина строки ненулевая, то:
		fieldsSlice := strings.Split(fields, ",") // Преобразуем строку, введённую пользователем (номера столбцов через запятую) в слайс, разделяя эту строку по сепаратору - запятой
		for _, f := range fieldsSlice { // Дальше мы пробегаемся по слайсу и каждый элемент (номер столбца) преобразуем к целочисленному типу...
			v, err := strconv.Atoi(f) // Преобразование к int может завершиться неуспешно в 2-х случаях:
			if err != nil {           // 1. Сама ф-ция Atoi завершилась с ошибкой, тогда возвращаем ошибку, сообщающую о невозможности распарсить переданную строку
				return errUnparsedField
			}
			if v <= 0 {				  // 2. Либо пользователем было передано отрицательное число в качестве номера столбца
				return errNegativeFieldValue
			}
			resultedFields = append(resultedFields, v-1) // ...затем добавляем этот элемент в результирующий слайс
		}
		o.fields = resultedFields // Присваиваем результирующий слайс полю структуры Options
		return nil
	}
}

func SetDelimeterOption(delim string) Option {
	return func(o *Options) error { // Если флаг "delimiter" не будет указан пользователем, параметр delim примет значение по умолчанию,
		o.delimeter = delim			// и в нашей структуре, инициализированной по умолчанию поле delimiter будет перезаписано тем же значением по умолчанию
		return nil
	}
}

func SetSeparatedOption(flag bool) Option { // Аналогично
	return func(o *Options) error {
		o.separated = flag
		return nil
	}
}
