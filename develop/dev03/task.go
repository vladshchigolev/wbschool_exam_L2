/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Утилита sort сортирует строки текстового файла
// По умолчанию сортировка осуществляется по всей строке,
// при желании можно разбить каждую строку по разделителю на столбцы и отсортировать строки по выбранному столбцу
func main() {
	sortByColumn := flag.Int("k", -1, "Указание колонки для сортировки") // По умолчанию сортируем по всей строке
	sortByNumbers := flag.Bool("n", false, "Сортировать по числовому значению")
	sortReverse := flag.Bool("r", false, "Сортировать в обратном порядке")
	noDuplicates := flag.Bool("u", false, "Не выводить повторяющиеся строки")
	ignoreEndSpace := flag.Bool("b", false, "Игнорировать хвостовые пробелы")
	checkSort := flag.Bool("c", false, "Проверять отсортированы ли данные")
	flag.Parse()

	fileName := os.Args[len(os.Args)-1] // Имя файла будет последним аргументом
	path, err := os.Getwd()             // Получаем текущую рабочую директорию
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadFile(path + "/" + fileName) // Читаем содержимое файла
	if err != nil {
		log.Fatal(err)
	}
	dataStr := string(data)
	result := Sort(dataStr, *sortByColumn, *sortByNumbers, *sortReverse, *noDuplicates, *ignoreEndSpace, *checkSort)

	print(result)
}

// Sort принимает в качестве аргументов: содержимое файла в виде строкового значения, набор "конфигурационных" переменных,
// значения которых установлены с помощью флагов
func Sort(data string, sortByColumn int, sortByNumbers, sortReverse, noDuplicates, ignoreEndSpace, checkSort bool) []string {
	if checkSort { // Если checkSort == true, проверяем, отсортированы ли данные
		result := strings.Split(data, "\n")    // Разделяем исходное строковое значение по "\n", наполняем слайс строк
		sorted := strings.Split(data, "\n")    // Делаем то же самое, только это слайс мы отсортируем
		sort.Strings(sorted)                   // Один из одинаковых слайсов отсортируем
		if reflect.DeepEqual(result, sorted) { // Сравним отсортированный с таким же, но несортированным, и если они равны, то:
			fmt.Println("sorted")
		} else { // Иначе:
			fmt.Println("unsorted")
		}
		return nil
	}
	if ignoreEndSpace {
		data = strings.TrimSpace(data) // Режем хвостовые пробелы
	}
	result := strings.Split(data, "\n") // Разделяем содержимое файла по "\n", на отдельные строки (т. к. именно их мы и будем сортировать), наполняем ими слайс
	if sortByColumn >= -1 {
		result = columnSort(result, sortByColumn)
	}

	if noDuplicates { // Не выводить повторяющиеся строки
		un := map[string]struct{}{}
		for _, r := range result {
			un[r] = struct{}{}
		}
		result = make([]string, 0, len(un))
		for k := range un {
			result = append(result, k)
		}
	}
	if sortByNumbers {
		nums := make([]int, 0, len(result)) // Создаём слайс целых чисел с cap равной количеству строк для сортировки
		if checkSort {                      // Проверяет, отсортировано ли содержимое файла
			if !sort.IntsAreSorted(nums) {
				return []string{"Need to sort"}
			} else {
				return []string{"No need to sort"}
			}
		}
		for k := range result { // Перебираем индексы строк
			n, err := strconv.Atoi(result[k])
			if err != nil {
				fmt.Println("error sorting by number:", err)
				os.Exit(1)
			}
			nums = append(nums, n)
		}
		if sortReverse {
			sort.Sort(sort.Reverse(sort.IntSlice(nums)))
		} else {
			sort.Ints(nums)
		}

		result = make([]string, 0, len(nums))
		for _, n := range nums {

			result = append(result, strconv.Itoa(n))
		}
	}

	if sortReverse {
		SortReverse(result)
	}

	return result
}

// columnSort сортирует строки по какому-то столбцу. Если в строке меньше столбцов она сдвигается
// вниз и сортируется по последнему слову в строке
func columnSort(data []string, ColumnNum int) []string {
	if ColumnNum == -1 { // Номер столбца по которому сортируем. По умолчанию - вся строка (до перевода строки)
		sort.Strings(data)
		return data
	}

	sort.Slice(data, func(i, j int) bool {
		lhs := strings.Split(data[i], " ")
		rhs := strings.Split(data[j], " ")
		if len(lhs) <= ColumnNum || len(rhs) <= ColumnNum { // Если количество столбцов в двух соседних строках меньше, чем н
			return lhs[0] < rhs[0] // Сортируем по возрастанию слова в столбце
		}
		return strings.Split(data[i], " ")[ColumnNum] <
			strings.Split(data[j], " ")[ColumnNum]
	})
	return data

}

// SortByNumbers сортирует строки
func SortByNumbers(text []string, row int) []string {
	var count int
	for i, v := range text { // сдвигаем вниз все строки, которые не числа
		words := strings.Split(v, " ")
		if len(words) < row {
			temp := text[i]
			for k := i; k < len(text)-1; k++ {
				text[k] = text[k+1]
			}
			text[len(text)-1] = temp
		} else {
			count++ // счетчик строк, которые числа
		}
	}
	return nil
}

func SortReverse(data []string) []string {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	return data
}

func print(text []string) {
	for _, v := range text {
		fmt.Println(v)
	}
}
