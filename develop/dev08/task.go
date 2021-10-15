/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

package main

import (
	"dev08/command"
	"dev08/shell"
	"flag"
	"log"
	"os"
)

var help bool

func init() {
	flag.BoolVar(&help, "help", false, "Показать помощь и выйти.")
}

func usage() {
	log.Printf(`ПРОСТЕЙШАЯ ОБОЛОЧКА ИНТЕРПРЕТАТОРА КОМАНДНОЙ СТОКИ. ДЛЯ ВЫХОДА ИСПОЛЬЗУЙТЕ \exit.`) // для регулярных выражений лучше использовать необработанные строки (raw strings, строки без интерпретации экранированных литералов)
	flag.PrintDefaults()
}

func showUsageAndExit(exitCode int) {
	usage()
	os.Exit(exitCode)
}

func main() {
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if help {
		showUsageAndExit(0)
	}

	s := shell.New()
	commander := command.New(s, os.Stdout, os.Stdin)
	if err := commander.Start(); err != nil {
		log.Fatal(err)
	}

}