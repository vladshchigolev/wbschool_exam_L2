package main

import "fmt"

// Отправитель
type button struct {
	command command // Интерфейс команды
}

func (b *button) press() {
	b.command.execute()
}

// Классы (типы) команд можно объединить под общим интерфейсом c единственным методом запуска.
// После этого одни и те же отправители смогут работать с различными командами, не привязываясь к их классам (типам).
type command interface {
	execute()
}

// Конкретная команда "ON"
type onCommand struct {
	device device // Интерфейс получателя
}

func (c *onCommand) execute() {
	c.device.on()
}

// Конкретная команда "OFF"
type offCommand struct {
	device device
}

func (c *offCommand) execute() {
	c.device.off()
}

// Интерфейс получателя
type device interface {
	on()
	off()
}

// Конкретный получатель
type tv struct {
	isRunning bool
}

func (t *tv) on() {
	t.isRunning = true
	fmt.Println("Turning tv on")
}

func (t *tv) off() {
	t.isRunning = false
	fmt.Println("Turning tv off")
}

// Клиентский код
func main() {
	tv := &tv{} // Создаём конкретного получателя

	onCommand := &onCommand{ // Создаём объект-команду, передаём ему ссылку на получателя
		device: tv,
	}

	offCommand := &offCommand{
		device: tv,
	}

	onButtonTV := &button{ // Создадим кнопку включения на самом телевизоре
		command: onCommand,
	}
	onButtonTV.press()

	offButtonTV := &button{ // Создадим кнопку выключения на самом телевизоре
		command: offCommand,
	}
	offButtonTV.press()

	onButtonRemote := &button{ // Создадим кнопку включения на пульте
		command: onCommand,
	}
	onButtonRemote.press()

	offButtonRemote := &button{ // Создадим кнопку выключения на пульте
		command: offCommand,
	}
	offButtonRemote.press()
	// Без объектов-команд пришлось бы дважды определять одинаковое поведение для кнопок на пульте и на самом телевизоре
}