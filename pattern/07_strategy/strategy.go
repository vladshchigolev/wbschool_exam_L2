// Cтратегия позволяет варьировать поведение исходного объекта (контекста) во время выполнения программы,
// подставляя в него различные объекты-поведения (например, отличающиеся балансом скорости и потребления ресурсов).
// Стратегия описывает разные способы произвести одно и то же действие,
// позволяя взаимозаменять эти способы в каком-то объекте контекста.

// Допустим, мы разрабатываем «In-Memory-Cache». Поскольку он находится внутри памяти, его размер ограничен.
// Как только он полностью заполнится, какие-то записи придется убрать для освобождения места. Эту функцию
// можно реализовать с помощью нескольких алгоритмов, самые популярные среди них:
// - Наиболее давно использовавшиеся (Least Recently Used – LRU): убирает запись, которая использовалась наиболее давно.
// - «Первым пришел, первым ушел» (First In, First Out — FIFO): убирает запись, которая была создана раньше остальных
// - Наименее часто использовавшиеся (Least Frequently Used — LFU): убирает запись, которая использовалась наименее часто.
// Проблема заключается в том, чтобы отделить кэш от этих алгоритмов для возможности их замены «на ходу». Помимо этого,
// класс кэша не должен меняться при добавлении нового алгоритма. В такой ситуации нам поможет паттерн "Стратегия".
package main

import "fmt"

// Все классы применяют одинаковый интерфейс, что делает алгоритмы взаимозаменяемыми внутри семейства.
type evictionAlgo interface {
	evict(c *cache)
}
// Конкретная стратегия 1
type fifo struct {
}

func (l *fifo) evict(c *cache) {
	fmt.Println("Evicting by fifo strategy")
}
// Конкретная стратегия 2
type lru struct {
}

func (l *lru) evict(c *cache) {
	fmt.Println("Evicting by lru strategy")
}
// Конкретная стратегия 3
type lfu struct {
}

func (l *lfu) evict(c *cache) {
	fmt.Println("Evicting by lfu strategy")
}
// Контекст
type cache struct {
	storage      map[string]string
	evictionAlgo evictionAlgo
	capacity     int
	maxCapacity  int
}

func initCache(e evictionAlgo) *cache { // initCache передаётся один из объектов-стратегий
	storage := make(map[string]string)
	return &cache{
		storage:      storage,
		evictionAlgo: e,
		capacity:     0,
		maxCapacity:  2,
	}
}

func (c *cache) setEvictionAlgo(e evictionAlgo) { // Метод для смены объекта-стратегии
	c.evictionAlgo = e
}

func (c *cache) add(key, value string) { // Перед каждым добавлением в кэш проверяем, не переполнен ли он
	if c.capacity == c.maxCapacity {
		c.evict() // Если кэш переполнен, вызываем метод evict(), который, в свою очередь, вызывает метод evict() установленной стратегии
	}
	c.capacity++
	c.storage[key] = value
}

func (c *cache) get(key string) {
	delete(c.storage, key)
}

func (c *cache) evict() { // Метод контекста evict() переадресует к методу evict() объекта-стратегии
	c.evictionAlgo.evict(c)
	c.capacity--
}

func main() {
	lfu := &lfu{} // Инициализируем стратегию
	cache := initCache(lfu) // Инициализируем объект-контекст (кэш) с какой-нибудь установленной стратегией

	cache.add("a", "1")
	cache.add("b", "2")

	cache.add("c", "3")

	lru := &lru{}
	cache.setEvictionAlgo(lru) // Поменяли стратегию

	cache.add("d", "4")

	fifo := &fifo{}
	cache.setEvictionAlgo(fifo)

	cache.add("e", "5")

}