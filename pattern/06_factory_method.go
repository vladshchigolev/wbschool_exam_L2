// ! В Go невозможно реализовать классический вариант паттерна Фабричный метод.
// ! Несмотря на это, все же можно реализовать базовую версию этого паттерна — Простая фабрика.
// -----------------------------------------------------------------------------------------------
// Когда весь наш код работает с объектами определённого класса (грузовиков, например), и бОльшая часть существующего кода жёстко привязана к этому классу или его подклассам,
// становится чрезвычайно сложным добавление в программу классов того же уровня (самолётов или поездов) без перелопачивания всей программы.
// Паттерн "Фабричный метод" предлагает создавать объекты не напрямую, используя конструктор конкретного объекта, например, а через вызов особого фабричного метода.

// Мы будем создавать разные типы музыкального инструмента фортепиано при помощи фабрики.
// Сперва мы создадим интерфейс iPiano, который определяет все методы будущих инструментов.
// Также имеем тип piano, который реализует интерфейс iPiano.
// Два конкретных инструмента — steinwayD274 и yamahaP45 — оба включают в себя структуру Piano и косвенно реализуют iPiano.
package main

import (
	"errors"
	"fmt"
)
//----------------------------------
// Сперва мы создадим интерфейс iPiano, который определяет общие методы, которые будут у всех будущих инструментов.
type iPiano interface {
	setType(newPianoType string)
	setColor(newColor string)
	getType() string
	getColor() string
}
//-----------------------------------
// Также имеем тип piano, который реализует интерфейс iPiano.
type Piano struct {
	pianoType  string
	color string
}

func (p *Piano) setType(newPianoType string) {
	p.pianoType = newPianoType
}

func (p *Piano) getType() string {
	return p.pianoType
}

func (p *Piano) setColor(newColor string) {
	p.color = newColor
}

func (p *Piano) getColor() string {
	return p.color
}
//---------------------------------------------------------------------------------------
// Два конкретных инструмента — steinwayD274 и yamahaP45 — оба включают в себя структуру Piano и косвенно реализуют iPiano.
type steinwayD274 struct {
	Piano
}

func newSteinwayD274(pianoType, color string) iPiano {
	return &steinwayD274{Piano: Piano{pianoType:  pianoType, color: color}}
}

func (sw *steinwayD274) String() string {
	return fmt.Sprintf("This is %s %s piano\n", sw.Piano.color, sw.Piano.pianoType)
}
//----------------------------------------------------------------------------------------

type yamahaP45 struct {
	Piano
}

func newYamahaP45(pianoType, color string) iPiano {
	return &yamahaP45{Piano: Piano{pianoType:  pianoType, color: color}}
}

func (yam *yamahaP45) String() string {
	return fmt.Sprintf("This is %s %s piano\n", yam.Piano.color, yam.Piano.pianoType)
}
//----------------------------------------------------------------------------------------------
// GetPiano служит фабрикой, которая создает инструмент нужного типа в зависимости от аргумента на входе. Клиентом служит main.
// Вместо прямого взаимодействия с объектами steinwayD274 или yamahaP45, она создает экземпляры различных инструментов,
// используя для контроля изготовления только параметры в виде строк.
func GetPiano(model, color string) (iPiano, error) {
	switch model {
	case "Steinway_D274":
		return newSteinwayD274("Acoustic", color), nil
	case "Yamaha_P45":
		return newYamahaP45("Electric", color), nil
	default:
		return nil, errors.New("wrong piano model passed")
	}
}
//--------------------------------
// Клиентский код
func main() {
	mySteinway, _ := GetPiano("Steinway_D274", "black")
	myYamaha, _ := GetPiano("Yamaha_P45", "white")

	printDetails(mySteinway)
	printDetails(myYamaha)
}

func printDetails(p iPiano) {
	fmt.Printf("Type: %s", p.getType())
	fmt.Println()
	fmt.Printf("Power: %d", p.getColor())
	fmt.Println()
}