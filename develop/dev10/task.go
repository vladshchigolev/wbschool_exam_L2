/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему серверу, программа должна завершаться через timeout.
*/

package main

import (
	"dev10/client"
	"errors"
	"flag"
	"log"
	"time"
)

var (
	errNotEnoughArguments = errors.New("telnet: укажите хост и порт для подключения")
	version               bool
	timeout               time.Duration // При подключении к несуществующему серверу, программа должна завершаться через timeout
)

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "Таймаут для соединения")
}

func usage() {
	log.Printf(`TELNET-КЛИЕНТ (С использованием TCP). 
ИСПОЛЬЗОВАНИЕ : ./telnet [COMMAND] <host> <port>
ДОСТУПНЫЕ КОМАНДЫ:`)
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	flag.Usage = usage // если запустить программу с флагом -h, будет вызвана usage
	flag.Parse()

	// Слайс Args содержит не флаговые аргументы командной строки
	args := flag.Args()
	if len(args) != 2 {
		log.Fatal(errNotEnoughArguments)
	}

	host, port := args[0], args[1]
	runner := client.New(host, port, timeout)
	if err := runner.Start(); err != nil {
		log.Fatal(err)
	}
}
