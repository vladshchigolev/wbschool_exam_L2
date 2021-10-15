// Цепочка обязанностей позволяет запускать обработчиков последовательно один за другим в том порядке,
// в котором они находятся в цепочке.
package main

import "fmt"
// Интерфейс обработчика
type department interface {
	execute(*patient)
	setNext(department)
}
// Конкретный обработчик
type reception struct {
	next department
}

func (r *reception) execute(p *patient) {
	if p.registrationDone {
		fmt.Println("Patient registration already done")
		r.next.execute(p)
		return
	}
	fmt.Println("Reception registering patient")
	p.registrationDone = true
	r.next.execute(p)
}

func (r *reception) setNext(next department) { // Параметр next типа интерфейса, т. к. следующим можно передать любое значение, имеющее методы execute() и setNext()
	r.next = next
}
// Конкретный обработчик
type doctor struct {
	next department
}

func (d *doctor) execute(p *patient) {
	if p.doctorCheckUpDone {
		fmt.Println("Doctor checkup already done")
		d.next.execute(p)
		return
	}
	fmt.Println("Doctor checking patient")
	p.doctorCheckUpDone = true
	d.next.execute(p)
}

func (d *doctor) setNext(next department) {
	d.next = next
}
// Конкретный обработчик
type cashier struct {
	next department
}
// Обработчик cashier является последним в цепи
func (c *cashier) execute(p *patient) {
	if p.paymentDone {
		fmt.Println("Payment Done")
	}
	fmt.Println("Cashier getting money from patient")
}

func (c *cashier) setNext(next department) {
	c.next = next
}

type patient struct {
	name              string
	registrationDone  bool
	doctorCheckUpDone bool
	medicineDone      bool
	paymentDone       bool
}

func main() {
	// Инициализируем объекты-обработчики в порядке, обратном тому, в котором через них должен пройти клиент
	// Также каждому обработчику вызовом соответствующего метода setNext() передаём указатель на другой обработчик,
	// метод execute() которого будет вызван дальше. Таким образом, связывая обработчики, мы формируем цепочку.
	// В любой момент можно вмешаться в существующую цепочку и переназначить связи так, чтобы убрать или добавить новое звено.
	cashier := &cashier{}

	//Set next for doctor department
	doctor := &doctor{}
	doctor.setNext(cashier)

	//Set next for reception department
	reception := &reception{}
	reception.setNext(doctor)

	patient := &patient{name: "Ivan"}
	//Patient visiting
	reception.execute(patient)
}