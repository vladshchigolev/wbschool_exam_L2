package client

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	normalExitMsg = "telnet: осуществляется выход...."
	//endOfFileExitMsg = "telnet: recieved EOF"
)

// Runner ...
type Runner struct {
	host    string        // Имя хоста для подключения (ip/доменное имя)
	port    string        // Номер порта
	timeout time.Duration // Таймаут для соединения
}

// Конструктор New создаёт значение типа (экземпляр класса) Runner и возвращает указатель на него
func New(host, port string, timeout time.Duration) *Runner {
	return &Runner{
		host:    host,
		port:    port,
		timeout: timeout,
	}
}

// Метод Start
func (r *Runner) Start() error {
	// Создаём telnet-клиент
	client := NewClient(r.host+":"+r.port, r.timeout, os.Stdin, os.Stdout) // После подключения STDIN программы должен записываться в сокет,
	// а данные, полученные из сокета, должны выводиться в STDOUT => в качестве значений io.Reader/Writer передаём os.Stdin/Stdout
	if err := client.BuildConnection(); err != nil { // Установление соединения
		return err
	}
	defer client.Close() // Отложенное освобождение подключения

	signalCh := make(chan os.Signal, 1)
	errorCh := make(chan error, 1) // По этому каналу будут отправляться/получаться значения ошибок
	// Cигналы следует ловить вручную, «подписавшись» на нужные типы сигналов:
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM) // Ctrl+C или kill <pid>

	go get(client, errorCh) // запустятся в новых горутинах
	go send(client, errorCh)

	wg := new(sync.WaitGroup) // Используем WaitGroup, чтобы main-горутина ждала завершения выполнения горутины на 57-ой строке
	wg.Add(1)
	// Ф-я ниже запускается в новой горутине. Она отвечает за завершение работы программы:
	// 1. Во-первых, она "держит" main-горутину, не позволяя ей завершиться раньше времени
	// 2. Она слушает события из 2-х каналов: signalCh, errorCh. Предполагается, что чтение произойдёт единожды из одного из каналов,
	// после чего выполнится оператор return и main-горутина разблокируется.
	go func() {
		defer wg.Done() // Уменьшим счётчик на 1 при завершении функции
		for {           // На каждой итерации цикла for в операторе select выполнится тот case, для которого будет готова операция чтения из канала
			select {
			// Сигнал о завершении программы может быть послан 2-х видов:
			case <-signalCh: // Если значение можно прочитать из этого канала, выводим сообщение, затем данная горутина завершается
				log.Println(normalExitMsg)
				return
			case err := <-errorCh:
				if err != nil {
					//if err == io.EOF { // Если ошибка связана с тем, что нечего больше читать,
					log.Println(err)
					//}
					return
				}
			default: // Если нечего прочитать из обоих каналов, переходим к следующей итерации for
				continue
			}
		}
	}()

	wg.Wait()
	return nil
}

// get и send завершатся самостоятельно в том случае, если не удалось по какой-то причине получить/отправить данные.
// В этом случае они отправят значение ошибки в канал, после чего вся программа должна завершиться
func send(c *Client, errorCh chan error) {
	for {
		if err := c.Send(); err != nil { // Если метод Send() вернул ошибку,
			errorCh <- err // отправим её в канал, что повлечёт завершение программы
			return
		}
	}
}

func get(c *Client, errorCh chan error) {
	for {
		if err := c.Get(); err != nil {
			errorCh <- err
			return
		}
	}
}
