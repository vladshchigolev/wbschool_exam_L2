package unpack

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	errInvalidString = errors.New("invalid string literal") // Создадим значение ошибки
)

const (
	backslashCode = 92 // Код в Unicode, представляющий "\"
)

// Unpack осуществляет примитивную распаковку строки
func Unpack(sample string) (string, error) { // вход: строка для распаковки, выход: распакованная строка, значение ошибки
	var res strings.Builder // Тип Builder используется для "сборки" строкового значения. Предоставляет набор методов. Например, метод Write(p []byte) добавляет содержимое p в буфер значения типа Builder
	runeSlice := []rune(sample) // Преобразуем исходную строку к слайсу рун
	for i := 0; i < len(runeSlice); i++ { // Перебираем индексы runeSlice (бежим по всем элементам runeSlice)
		// Т. к. значение руны - целое число (код символа в юникоде, rune то же самое, что и int32).
		if runeSlice[i] == backslashCode { // Если i-ый элемент - "\" ...
			i++ // ...то в буфер билдера помещаем следующий за "\" элемент. Сейчас мы увеличили i на единицу, но в самом конце этой итерации выполнится ещё одна инструкция i++ (та, что в заголовке цикла)
			res.WriteRune(runeSlice[i])
		} else if n, err := strconv.Atoi(string(runeSlice[i])); err == nil { // Если удалось успешно преобразовать string к int (значение ошибки != nil), попадаем внутрь if
			// Иначе если i-ый элемент - цифра, может быть 2 варианта:
			// 1. Когда входная строка - некорректная (возвращаем пустую строку и значение ошибки). Некорректной она может быть в следующих случаях:
			if i == 0 || // 1.1 Если i-тый элемент - цифра, являющаяся первым элементом в исходной строке
				(i > 0 && unicode.IsDigit(runeSlice[i-1]) && // 1.2 ИЛИ такое сочетание условий:
					(i > 1 && runeSlice[i-2] != backslashCode)) { // если предыдущий символ - цифра, И эта цифра НЕ ЭКРАНИРОВАНА (пред-предыдущий символ - "\"), иными словами, идут 2 цифры подряд
				// Условия "i > 0" и "i > 1" говорят: "при условии, что эти предыдущие символы вообще существуют"
				return "", errInvalidString
			}
			// 2. Если вышеописанные условия НЕ ВЫПОЛНЕНЫ, значит предшествующий нашей цифре символ уже был ЕДИНОЖДЫ записан в буфер, записываем его ещё n-1 раз
			// если i-ый элемент цифра и все ок
			// записываем n-1 раз предыдущий элемент
			// n-1 т.к. исходник уже записан
			res.WriteString(strings.Repeat(string(runeSlice[i-1]), n-1))
		} else {
			res.WriteRune(runeSlice[i]) // Если этот элемент НЕ "\" и НЕ цифра, добавляем его ЕДИНОЖДЫ в буфер Builder, потом если за ним в исходной строке пойдет цифра (например n),
			// мы добавим наш элемент в буфер ещё n-1 раз (43-я строка)
		}
	}
	return res.String(), nil
}