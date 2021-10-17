/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

package main

import (
	"dev06/cut"
	"dev06/managers"
	"flag"
	"log"
	"os"
)

var (
	fields      string
	delimeter   string
	isSeparated bool
	help        bool
)

func init() { // init() запускается сразу же после импорта пакета, используется при необходимости инициализации приложения в определенном состоянии
	flag.StringVar(&fields, "f", "", "Выбрать поля (колонки). Целые положительные числа, через запятую в кавычках")
	flag.StringVar(&delimeter, "d", "\t", "Использовать другой разделитель")
	flag.BoolVar(&isSeparated, "s", false, "Выводить только строки с разделителем")
	flag.BoolVar(&help, "help", false, "Показать помощь и выйти.")
}

func usage() {
	log.Printf(`Уитилита для обрезки строк. Использование: ./[название_исполняемого_файла] [опции]
Доступные опции:`)
	flag.PrintDefaults() // Выводит доступные флаги, их usage и дефолтное значение
}

func showUsageAndExit(exitCode int) {
	usage()
	os.Exit(exitCode)
}

func main() {
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse() // Вызываем Parse чтобы распарсить наши флаги и записать их в нужные переменные
	//usage()
	if help { // Флаг принимает значение булевого типа
		showUsageAndExit(0)
	}

	conManager := managers.NewConsoleManager(os.Stdin, os.Stdout) // В нашем случае читать будем из stdin.
	// Stdin, Stdout - это открытые файлы, указывающие на файловые дескрипторы стандартных ввода и вывода
	options := []cut.Option{ // Здесь с помощью литерала слайса создаём слайс значений Option (функций с сигнатурой "func(o *Options) error")
		cut.SetFieldsOption(fields),         // Здесь каждая
		cut.SetDelimeterOption(delimeter),   // инструкция вернёт
		cut.SetSeparatedOption(isSeparated), // по функции (по значению типа Option)
		// Аргументы, передаваемые этим трём функциям, взяты из пользовательского ввода ()
		// Если один/несколько/все из этих аргументов - значения по умолчанию для флагов, то:
		// 1. ф-ция, возвращённая ф-цией SetFieldsOption, если дать ей на вход "", вернёт ошибку
	}

	newCut := cut.New(conManager).ApplyOptions(options...) // Создаём новый экземпляр ManCut, передаем ему менеджера, применяем нужные опции (вызываем все функции из слайса)
	if err := newCut.Cut(); err != nil {                   // Здесь собственно и вызываем метод "Cut" типа ManCut
		log.Fatal(err)
	}

	if err := newCut.OutputResult(); err != nil {
		log.Fatal(err)
	}
}
