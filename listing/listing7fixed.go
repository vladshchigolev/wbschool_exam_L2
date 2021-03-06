package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)
// Что происходит: в main-горутине дважды вызывается функция, в которой:
// 1. создается небуферизованный канал
// 2. в новой горутине вызывается функция, перебирающая в цикле элементы массива и отправляющая их значения по созданному каналу
func asChan(vs ...int) <-chan int { // Функция возвращает канал, по которому она отправляет значения
	// На вызывающей стороне из канала можно ТОЛЬКО ЧИТАТЬ
	// Зачем вообще возвращать канал? - Чтобы писать в канал или читать из него, нужно, во-первых, иметь к нему доступ (ссылку на структуру)
	// длина массива vs определяется количеством инициализаторов
	c := make(chan int) // c - локальная переменная
	//fmt.Println(reflect.TypeOf(vs))
	go func() {
		for _, v := range vs { // Перебираем элементы массива
			c <- v // Выполняем операцию отправки по каналу (данная горутина заблокируется, пока значение не будет получено)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // Вызываем time.Sleep
		}

		close(c) // после успешного отправления (буфер = 0) закрываем канал, чтобы получающая горутина (если получение огранизовано в цикле) не заблокировалась навсегда
	}()
	return c // Возвращаем канал
}
// В какой-то момент горутины, созданные в функциях asChan отправив первые значения (первую пару: каждая горутина по одному значению) по своим каналам заблокируются,
// ожидая получения значения кем-нибудь (операции чтения из этих каналов производит горутина, созданная в функции newMerge)
// Здесь непосредственно происходит сливание 2-х каналов в один
func newMerge(a, b <-chan int) <-chan int { // вход: 2 канала, выход: 1 общий
	var wg sync.WaitGroup // Чтобы узнать, когда последняя горутина закончит работу
						  // (эта горутина может не быть последней из запущенных на выполнение),
						  // нужно увеличивать счетчик перед запуском каждой горутины (wg.Add()) и уменьшать
						  // его после завершения (wg.Done()). Для этого требуется счетчик особого рода - sync.WaitGroup, с которым могут
						  // безопасно работать несколько горутин и который предоставляет возможность ожидания (wg.Wait())
	 					  // пока он не станет равным нулю.
	out := make(chan int) // В этот канал сливаем

	output := func(c <-chan int) { // вход функции: канал, из которого будем читать...
		for data := range c {	   // выйдем из цикла, когда канал на той стороне закроется
			out <- data			   // ...и писать значение в канал, предназначенный для сливания
		}
		wg.Done()				   // уменьшаем счетчик на 1
	}
	wg.Add(2) // увеличиваем счетчик на 2 (ниже следует запуск 2-х горутин)
	go output(a) // Будем одновременно читать из обоих каналов и писать в общий, не надо выбирать из какого канала сейчас прочитать значение
	go output(b)

	go func() { // Почему в отдельной горутине? Нам нужно идти дальше по main-горутине (читать значения из общего канала out и выводить их).
		// Если мы выполним этот код в main-горутине, будет deadlock, т. к. она будет ждать, когда те две закончат работу, чтобы закрыть канал,
		// а они в свою очередь будут ждать, когда кто-нибудь прочитает первые отправленные ими значения по каналу (это должна быть main-горутина)
		wg.Wait() // ждём, пока значение счётчика не станет равным 0
		close(out) // после чего закрываем канал
	}()
	return out
}

func main() {
	// Здесь вызовы asChan и merge синхронны (конец выполения 1-ой совпадает с началом 2-ой)
	a := asChan(1, 3, 5, 7) // asChan вернула канал
	b := asChan(2, 4 ,6, 8) // asChan вернула канал
	c := newMerge(a, b) // Передаем newMerge каналы, возвращённые в результате 2-х вызовов asChan
	for v := range c { // производим чтение из канала. выйдем из цикла, когда канал на той стороне закроется
		fmt.Println(v)
	}
}