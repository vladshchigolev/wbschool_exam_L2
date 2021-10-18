![Builder pattern](https://refactoring.guru/images/patterns/diagrams/builder/problem1.png)

Создав кучу подклассов для всех конфигураций объектов, вы можете излишне усложнить программу.

![Builder pattern](https://refactoring.guru/images/patterns/diagrams/builder/problem2.png)

Конструктор со множеством параметров имеет свой недостаток: не все параметры нужны большую часть времени.

>Паттерн "Строитель" предлагает вынести конструирование объекта за пределы его собственного класса, поручив это дело отдельным объектам, называемым строителями.

![Builder pattern](https://refactoring.guru/images/patterns/diagrams/builder/solution1.png)

Строитель позволяет создавать сложные объекты пошагово. Промежуточный результат всегда остаётся защищён.