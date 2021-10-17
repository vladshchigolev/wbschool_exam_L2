// Основная идея в том, что программа может находиться в одном из нескольких состояний, которые всё время сменяют друг друга.
// Набор этих состояний, а также переходов между ними, предопределён и конечен. Находясь в разных состояниях,
// программа может по-разному реагировать на одни и те же события, которые происходят с ней.
// Паттерн "Состояние" нужно использовать в случаях, когда объект может иметь много различных состояний,
// которые он должен менять в зависимости от конкретного поступившего запроса.
// В нашем примере, торговый автомат может быть в одном из множества состояний, которые непрерывно меняются.
// Допустим, что автомат находится в режиме itemRequested. Как только произойдет действие «Ввести деньги»,
// он сразу же перейдет в состояние hasMoney. В зависимости от состояния торгового автомата, в котором он находится в данный момент,
// он может по-разному отвечать на одни и те же запросы. Например, если пользователь хочет купить предмет, машина выполнит действие,
// если она находится в режиме hasItemState, и отклонит запрос в режиме noItemState.
package main

import "fmt"

type vendingMachine struct {
	// Благодаря тому, что объекты состояний будут иметь общий интерфейс,
	// контекст сможет делегировать работу состоянию, не привязываясь к его классу (типу).
	// Поведение контекста можно будет изменить в любой момент, подключив к нему другой объект-состояние.
	hasItem       state // vendingMachine будет хранить ссылки на все объекты-состояния, но делегировать работу
	itemRequested state // будет тому из них, кто сейчас находится в currentState
	hasMoney      state
	noItem        state

	currentState state // В зависимости от значения currentState, одни и те же методы, вызванные у vendingMachine будут производить разные действия
	// Именно к объекту, который сейчас находится в currentState, vendingMachine будет перенаправлять запросы
	itemCount int
	itemPrice int
}
// newVendingMachine конструктор объекта vendingMachine,
func newVendingMachine(itemCount, itemPrice int) *vendingMachine {
	v := &vendingMachine{
		itemCount: itemCount,
		itemPrice: itemPrice,
	}
	// Инициализируем все состояния...
	hasItemState := &hasItemState{
		vendingMachine: v,
	}
	itemRequestedState := &itemRequestedState{
		vendingMachine: v,
	}
	hasMoneyState := &hasMoneyState{
		vendingMachine: v,
	}
	noItemState := &noItemState{
		vendingMachine: v,
	}
	// ...и присваиваем их полям контекста
	// Вместо того, чтобы хранить код всех состояний, первоначальный объект (v), называемый контекстом,
	// будет содержать ссылку на один из объектов-состояний и делегировать ему работу, зависящую от состояния.
	v.setState(hasItemState)
	v.hasItem = hasItemState
	v.itemRequested = itemRequestedState
	v.hasMoney = hasMoneyState
	v.noItem = noItemState
	return v
}
// Поскольку автомат при вызове какого-либо метода находится в одном из множества состояний,
// мы обращаемся к реализации соответствующего. В итоге результат вызова каждого метода будет
// отличаться в зависимости от currentState
func (v *vendingMachine) requestItem() error {
	return v.currentState.requestItem()
}

func (v *vendingMachine) addItem(count int) error {
	return v.currentState.addItem(count)
}

func (v *vendingMachine) insertMoney(money int) error {
	return v.currentState.insertMoney(money)
}

func (v *vendingMachine) dispenseItem() error {
	return v.currentState.dispenseItem()
}

func (v *vendingMachine) setState(s state) { // Объекты-состояния будут вызывать этот метод vendingMachine, после того как сами отработали,
	v.currentState = s						 // чтобы сменить состояние vendingMachine.
}

func (v *vendingMachine) incrementItemCount(count int) {
	fmt.Printf("Adding %d items\n", count)
	v.itemCount = v.itemCount + count
}
