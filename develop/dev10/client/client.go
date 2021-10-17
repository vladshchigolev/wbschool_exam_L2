package client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var (
	errBuildDialWithTimeout   = errors.New("telnet: не удалось установить соединение")
	errCloseConnection        = errors.New("telnet: не удалось завершить соединение")
	errRequestToBadConnection = errors.New("telnet: не удалось получить/отправить данные по сети")
)

// Тип (класс) Client представляет telnet-клиент
type Client struct {
	addr            string        // [имя_хоста/ip:номер_порта]
	timeout         time.Duration // Таймаут для соединения
	conn            net.Conn      // Представляет объект-соединение с удаленным хостом. Имеет методы Read/Write, => можно посылать/получать данные по сети
	inputDataReader *bufio.Reader // Чтение из STDIN
	connDataReader  *bufio.Reader // Реализует чтение данных из conn через буфер
	writer          io.Writer     // Вывод в STDOUT
}

// Конструктор NewClient создаёт значение типа (экземпляр класса) Client и возвращает указатель на него
func NewClient(addr string, timeout time.Duration, inputDataReader io.Reader, writer io.Writer) *Client {
	return &Client{
		addr:            addr,
		timeout:         timeout,
		inputDataReader: bufio.NewReader(inputDataReader),
		writer:          writer,
	}
}

// Метод BuildConnection устанавливает соединение с удалённым хостом
func (c *Client) BuildConnection() error {
	// DialTimeout подключается к указанной сети
	conn, err := net.DialTimeout("tcp", c.addr, c.timeout) // net.DialTimeout() применяется для отправки запросов к ресурсам в сети
	// вход: network - тип протокола, address - адрес ресурса, timeout - таймаут для соединения
	// Функция возвращает объект, который реализует интерфейс net.Conn. В нашем случае это объект (значение) типа TCPConn.
	// Так как net.Conn реализует интерфейсы io.Reader и io.Writer, то в данный объект (типа TCPConn) можно записывать данные - фактически посылать
	// по сети данные и можно считывать из него данные - получать данные из сети.
	if err != nil {
		return errBuildDialWithTimeout
	}
	c.conn = conn
	c.connDataReader = bufio.NewReader(conn) // Создание потока ввода через буфер
	log.Println("telnet: подключение успешно установлено")
	return nil
}

// Get считывает данные из подключения и пишет их в stdout. Get будет запущена в цикле, который в свою очередь запущен в отдельной горутине
func (c *Client) Get() error {
	text, err := c.connDataReader.ReadString('\n') // Поскольку идет построчное считывание, то каждая строка считывается из потока, пока не будет обнаружен символ перевода строки \n
	if err != nil {
		//if err == io.EOF {
		//	err = errRequestToBadConnection
		//}
		return err
	}
	if _, err := fmt.Fprint(c.writer, text); err != nil { // Пишем содержимое text (получено из соединения) в STDOUT
		return err
	}
	return nil
}

// Send ...
func (c *Client) Send() error {
	text, err := c.inputDataReader.ReadString('\n') // Аналогично методу Get(), только здесь читаем построчно из stdin...
	if err != nil {
		return err
	}
	if _, err := c.conn.Write([]byte(text)); err != nil { // ...и пишем это в соединение.
		return errRequestToBadConnection
	}
	return nil

}

// Close ...
func (c *Client) Close() error {
	if err := c.conn.Close(); err != nil {
		return errCloseConnection
	}
	return nil
}
