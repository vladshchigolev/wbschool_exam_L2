package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
var (
	after int
	before int
	contextText int
	countBool bool
	ignoreCase bool
	invert bool
	fixed bool
	lineNum bool
	filePath string
	help bool
	isNotDefault bool // Если не будет вызвана ни одна функция-обработчик, isNotDefault так и останется "false" и будет вызвана ф-ция Default, просто возвращающая только строки с совпадением
)
func main() {
	flag.IntVar(&after, "A", 0, "печатать N строк после совпадения")
	flag.IntVar(&before, "B", 0, "печатать N строк до совпадения")
	flag.IntVar(&contextText, "C", 0, "печатать ±N строк вокруг совпадения")
	var count int
	flag.BoolVar(&countBool, "c", false, "количество строк")
	flag.BoolVar(&ignoreCase, "i", false, "игнорировать регистр")
	flag.BoolVar(&invert, "v", false, "вместо совпадения, исключать")
	flag.BoolVar(&fixed, "F", false, "точное совпадение со строкой, не паттерн")
	flag.BoolVar(&lineNum, "n", false, "печатать номер строки")
	flag.StringVar(&filePath,"fp", "", "указать абсолютный путь до файла")
	flag.BoolVar(&help, "h", false, "показать помощь и выйти")
	flag.Parse() // Парсим флаги и записываем их значения в нужные переменные
	// В Golang у нас есть пакет под названием os, который содержит массив под названием «Args».
	// Args - это массив строк, содержащий все переданные аргументы командной строки (разделителем между аргументами выступает символ пробела).
	// Первым аргументом всегда будет имя программы. Последним мы условимся, что будет искомое строковое значение
	phrase := os.Args[len(os.Args)-1] // str будет содержать искомое фразу

	if help {
		printHelpAndExit(0)
	}

	t, err := ioutil.ReadFile(filePath) // Читает файл, возвращает последовательность байтов и значение ошибки
	if err != nil {
		log.Fatal(err)
	}
	text := string(t) // Приводим слайс байтов к строковому значению

	fmt.Println("===ИСХОДНЫЙ ТЕКСТ===")
	fmt.Println(text) // text содержит исходный текст. В дальнейшем её значение будет модифицировано (присвоен результат работы программы)
	fmt.Println("===РЕЗУЛЬТАТ===")

	if lineNum {
		fmt.Println("Find on", LineNum(phrase, text)+1, "line")
	}

	if countBool {
		count = Count(text)
		fmt.Println("Count of lines", count)
	}

	if ignoreCase {
		text = IgnoreCase(phrase, text, after, before, contextText)
	}

	if after > 0 { // Если пользователь указал с помощью флага значение переменной after, отличное от нуля, значит он хочет чтобы было напечатано "after" строк после совпадения
		text = After(phrase, text, after) // Вызываем ф-цю, возвращающую требуемый пользователем результат
	}
	// Аналогично вышеописанному
	if before > 0 {
		text = Before(phrase, text, before)
	}
	// Аналогично вышеописанному
	if contextText > 0 {
		text = ContextText(phrase, text, contextText)
	}

	if invert {
		text = Invert(phrase, text)
	}

	if fixed {
		fmt.Println(Fixed(phrase, text))
	}
	// Лишь в конце проверяем, нужно ли выполнить действие по умолчанию (если пользователь не указал ни одного значения флага)
	if !isNotDefault {
		text = Default(phrase, text)
	}
	fmt.Println(text)
}
func printHelpAndExit(exitCode int) {
	log.Printf(`Утилита grep. Использование: ./[название_исполняемого_файла] [опции] [искомый текст].
Для того, чтобы вывести ТОЛЬКО строки, в которых встречается искомая фраза/значение, используйте любой из флагов -A, -B, -C со значением "0"
Доступные опции:`)
	flag.PrintDefaults()
	os.Exit(exitCode)
}
// Default просто возвращает строки, содержащие совпадение, без применения различных опций
func Default(str, file string) (res string) {
	expr, err := regexp.Compile(str) // Ф-ция Compile возвращает регулярное выражение
	if err != nil {
		fmt.Println(err)
	}

	rows := strings.Split(file, "\n")
	for _, v := range rows {
		if expr.MatchString(v) {
			res += v
		}
	}
	return
}
// After возвращает искомую строку и count строк ПОСЛЕ нее
func After(str, file string, count int) (res string) { // вход: 1. Значение, строку в котором оно содержится, нужно найти. 2. Набор текстовых данных, в котором будем искать совпадение. 3. Кол-во строк после строки, где есть совпадение, для вывода
	//expr := regexp.MustCompile(str)
	expr, err := regexp.Compile(str) // Ф-ция Compile возвращает регулярное выражение
	if err != nil {
		fmt.Println(err)
	}
	rows := strings.Split(file, "\n") // Весь исходный набор текста разбиваем построчно на фрагменты и заполняем ими слайс строк
	for i, v := range rows { // Анализируем каждую строчку исходного текста
		if expr.MatchString(v) { // Метод MatchString проверяет, есть ли в строке вхождения регулярного выражения. Результат - true/false
			res += v + "\n" // Используем оператор конкатенации для того, чтобы добавить строку исходного текста (с переводом строки) с совпадением в РЕЗУЛЬТИРУЮЩУЮ строковую переменную res
			if !(i == len(rows)-1) { // Если строчка, в которой обнаружено совпадение - последняя, не добавляем n строчек после неё к результату, т. к. их больше нет
				for j, k := range rows { // Значение v (одну из строчек исходного текста, в которой обнаружено совпадение) сравниваем с КАЖДОЙ строчкой этого же текста (сравниваются ИНДЕКСЫ (положение в массиве))
					if j > i && j <= i+count { // Дальше отбираем строчки, находящиеся ПОСЛЕ строчки с совпадением, но в то же время не выходящие за границу (after, или лок. переменная count - то есть кол-во строк после совпадения)
						//if j == len(rows)-1 || j == i+count { // Если строчка прошла предыдущий барьер, то: если она последняя в тексте или последняя в "разрешенном диапазоне" вывода
						//	res += k // Добавляем ее к результату БЕЗ перевода строки (после неё уже точно никаких строк добавлено не будет)
						//} else {
						res += k + "\n" // Иначе добавляем с переводом строки
						//}
					}
				}
			}
		}
	}
	isNotDefault = true
	return
}

// Before возвращает искомую строку и count строк ДО нее (ф-ция, симметричная After)
func Before(str, file string, count int) (res string) {
	expr, err := regexp.Compile(str) // Ф-ция Compile возвращает регулярное выражение
	if err != nil {
		fmt.Println(err)
	}
	// Аналогично функции After
	rows := strings.Split(file, "\n")
	for i, v := range rows {
		if expr.MatchString(v) {
			for j, k := range rows {
				if j < i && j >= i-count {
					res += k + "\n"
				}
			}
			res += v // Здесь небольшое отличие от ф-ции After: строчку, в которой обнаружено совпадение, мы добавим к результату последней, чтобы было симметрично
		}
	}
	isNotDefault = true
	return
}

// ContextText возвращает искомую строку и count строк ВОКРУГ нее
func ContextText(str, file string, count int) (res string) {
	expr, err := regexp.Compile(str) // Ф-ция Compile возвращает регулярное выражение
	if err != nil {
		fmt.Println(err)
	}
	rows := strings.Split(file, "\n")
	for i, v := range rows {
		if expr.MatchString(v) {
			for j, k := range rows { // Цикл добавляет к результату
				if j < i && j >= i-count {
					res += k + "\n"
				}
			}
			res += v + "\n"
			for j, k := range rows {
				if j > i && j <= i+count {
					if j == len(rows)-1 || j == i+count {
						res += k
					} else {
						res += k + "\n"
					}
				}
			}
		}
	}
	isNotDefault = true
	return
}

// Count возвращает кол-во строк в исходном тексте
func Count(file string) (count int) {
	rows := strings.Split(file, "\n")
	for range rows {
		count++
	}
	isNotDefault = true
	return
}

// IgnoreCase возвращает искомую строку игнорируя регистр
func IgnoreCase(str string, file string, after int, before int, context int) (res string) {

	str = strings.ToLower(str) // Возвращает строку, где все символы str приведены к нижнему регистру
	file = strings.ToLower(file) // Сам текст, по которому будем искать, привести к нижнему регистру
	// Здесь без if-а должно быть действие по умолчанию
	if after > 0 {
		res = After(str, file, after)
		after = 0
	}
	if before > 0 {
		res = Before(str, file, before)
		before = 0
	}
	// Аналогично вышеописанному
	if context > 0 {
		res = ContextText(str, file, contextText)
		contextText = 0
	}
	return res
}

// Invert возвращает текст без искомой строки
func Invert(str, file string) (res string) {
	expr, err := regexp.Compile(str) // Ф-ция Compile возвращает регулярное выражение
	if err != nil {
		fmt.Println(err)
	}
	rows := strings.Split(file, "\n")
	for i, v := range rows {
		if expr.MatchString(v) {
			continue
		}
		if i == len(rows)-1 {
			res += v
		} else {
			res += v + "\n"
		}
	}
	isNotDefault = true
	return
}

// Fixed возвращает true, если имеется точное совпадение, false в противном случае
func Fixed(str, file string) bool {
	isNotDefault = true
	rows := strings.Split(file, "\n")
	for _, v := range rows {
		if v == str {
			return true
		}
	}
	return false
}

// LineNum возвращает номер исходной строки
func LineNum(str, file string) int {
	expr, err := regexp.Compile(str) // Ф-ция Compile возвращает регулярное выражение
	if err != nil {
		fmt.Println(err)
	}
	rows := strings.Split(file, "\n")
	for i, v := range rows {
		if expr.MatchString(v) {
			return i
		}
	}
	return -1
}