package main

import "fmt"

// Абстрактная фигура
type shape interface {
	getType() string
	accept(visitor)
}
// ------------------------
// Конкретная фигура (квадрат)
type square struct {
	side int
}

func (s *square) accept(v visitor) {
	v.visitForSquare(s)
}

func (s *square) getType() string {
	return "Square"
}
// ---------------------------------
// Конкретная фигура (круг)
type circle struct {
	radius int
}

func (c *circle) accept(v visitor) {
	v.visitForCircle(c)
}

func (c *circle) getType() string {
	return "Circle"
}
// ----------------------------------
// Конкретная фигура (треугольник)
type rectangle struct {
	l int
	b int
}

func (t *rectangle) accept(v visitor) {
	v.visitForRectangle(t)
}

func (t *rectangle) getType() string {
	return "rectangle"
}
// -------------------------------------
// Абстрактный посетитель (интерфейс посетителя)
type visitor interface {
	visitForSquare(*square)
	visitForCircle(*circle)
	visitForRectangle(*rectangle)
}
// Функции выше позволят нам добавлять функционал для квадратов, кругов и треугольников соответственно.
// При этом мы не можем оставить только один метод visit(shape) в интерфейсе посетителя, т. к. Go не поддерживает
// перегрузку методов, поэтому нельзя иметь методы с одинаковыми именами, но разными параметрами.
// ----------------------------------------
// Конкретный посетитель (вычисляет площадь)
// Поскольку конкретная фигура (квадрат, например) в качестве аргумента её метода accept() принимает значение интерфейсного
// типа visitor, для соответствия которому нужно обладать методами для всех конкретных фигур (ещё круг и треугольник),
// нужно реализовать аналогичное поведение (вычисление площади) этим посетителем для остальных фигур. Тогда посетитель,
// имея все нужные методы, без проблем может быть передан аргументом методу accept() любой из фигур.
type areaCalculatorVisitor struct {
	area int
}

func (a *areaCalculatorVisitor) visitForSquare(s *square) {
	// Вычисляет площадь квадрата.
	// Затем присваивает результат переменной (полю) area значения areaCalculatorVisitor.
	fmt.Println("Calculating area for square")
}

func (a *areaCalculatorVisitor) visitForCircle(s *circle) {
	fmt.Println("Calculating area for circle")
}
func (a *areaCalculatorVisitor) visitForRectangle(s *rectangle) {
	fmt.Println("Calculating area for rectangle")
}
// --------------------------------------------------------------
// Конкретный посетитель (вычисляет периметр)
type perimeterCalculatorVisitor struct {
	perimeter int
}

func (a *perimeterCalculatorVisitor) visitForSquare(s *square) {

	fmt.Println("Calculating middle point coordinates for square")
}

func (a *perimeterCalculatorVisitor) visitForCircle(c *circle) {
	fmt.Println("Calculating middle point coordinates for circle")
}
func (a *perimeterCalculatorVisitor) visitForRectangle(t *rectangle) {
	fmt.Println("Calculating middle point coordinates for rectangle")
}

func main() {
	square := &square{side: 2}
	circle := &circle{radius: 3}
	rectangle := &rectangle{l: 2, b: 3}

	areaCalculator := &areaCalculatorVisitor{}

	square.accept(areaCalculator)
	circle.accept(areaCalculator)
	rectangle.accept(areaCalculator)

	fmt.Println()
	perimeterCalculator := &perimeterCalculatorVisitor{}
	square.accept(perimeterCalculator)
	circle.accept(perimeterCalculator)
	rectangle.accept(perimeterCalculator)
}