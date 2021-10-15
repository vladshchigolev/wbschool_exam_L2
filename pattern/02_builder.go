// Паттерн "Строитель" предлагает вынести конструирование объекта за пределы его собственного класса,
// поручив это дело отдельным объектам, называемым строителями. Строитель позволяет создавать сложные объекты пошагово
// (например, построитьСтены, вставитьДвери и др.). Чтобы создать объект, нужно поочерёдно вызывать методы строителя.
// Причём не нужно запускать все шаги, а только те, что нужны для производства объекта определённой конфигурации.
package main

import "fmt"
// Зачастую один и тот же шаг строительства может отличаться для разных вариаций производимых объектов.
// Например, деревянный дом потребует строительства стен из дерева, а каменный — из камня.
// В этом случае можно создать несколько классов (типов) строителей, выполняющих одни и те же шаги, но по-разному.
// Интерфейс строителя объявляет шаги конструирования продуктов, общие для всех видов строителей.
type iBuilder interface {
	setWindowType()
	setDoorType()
	setNumFloor()
	getHouse() house
}

func getBuilder(builderType string) iBuilder {
	if builderType == "normal" {
		return &normalBuilder{}
	}

	if builderType == "igloo" {
		return &iglooBuilder{}
	}
	return nil
}
// Конкретные строители реализуют строительные шаги, каждый по-своему.
// Конкретные строители могут производить разнородные объекты, не имеющие общего интерфейса.
type normalBuilder struct {
	windowType string
	doorType   string
	floor      int
}

func newNormalBuilder() *normalBuilder {
	return &normalBuilder{}
}

func (b *normalBuilder) setWindowType() {
	b.windowType = "Wooden Window"
}

func (b *normalBuilder) setDoorType() {
	b.doorType = "Wooden Door"
}

func (b *normalBuilder) setNumFloor() {
	b.floor = 2
}

func (b *normalBuilder) getHouse() house {
	return house{
		doorType:   b.doorType,
		windowType: b.windowType,
		floor:      b.floor,
	}
}

type iglooBuilder struct {
	windowType string
	doorType   string
	floor      int
}

func newIglooBuilder() *iglooBuilder {
	return &iglooBuilder{}
}

func (b *iglooBuilder) setWindowType() {
	b.windowType = "Snow Window"
}

func (b *iglooBuilder) setDoorType() {
	b.doorType = "Snow Door"
}

func (b *iglooBuilder) setNumFloor() {
	b.floor = 1
}

func (b *iglooBuilder) getHouse() house {
	return house{
		doorType:   b.doorType,
		windowType: b.windowType,
		floor:      b.floor,
	}
}

type house struct {
	windowType string
	doorType   string
	floor      int
}
// Можно выделить вызовы методов строителя в отдельный класс, называемый директором.
// В этом случае директор будет задавать порядок шагов строительства, а строитель — выполнять их.
type director struct {
	builder iBuilder
}
// Для того, чтобы создать нового директора, нужно передать функции-конструктору какого-нибудь строителя
func newDirector(b iBuilder) *director {
	return &director{
		builder: b,
	}
}

func (d *director) setBuilder(b iBuilder) {
	d.builder = b
}

func (d *director) buildHouse() house {
	d.builder.setDoorType()
	d.builder.setWindowType()
	d.builder.setNumFloor()
	return d.builder.getHouse()
}

func main() {
	normalBuilder := getBuilder("normal")
	iglooBuilder := getBuilder("igloo")

	director := newDirector(normalBuilder) // normalBuilder строит обычный дом
	// Мы не напрямую вызываем методы строителя, а косвенно, через директора
	normalHouse := director.buildHouse()

	fmt.Printf("Normal House Door Type: %s\n", normalHouse.doorType)
	fmt.Printf("Normal House Window Type: %s\n", normalHouse.windowType)
	fmt.Printf("Normal House Num Floor: %d\n", normalHouse.floor)

	director.setBuilder(iglooBuilder)
	iglooHouse := director.buildHouse()

	fmt.Printf("\nIgloo House Door Type: %s\n", iglooHouse.doorType)
	fmt.Printf("Igloo House Window Type: %s\n", iglooHouse.windowType)
	fmt.Printf("Igloo House Num Floor: %d\n", iglooHouse.floor)

}
