package main

import (
	"fmt"
	"os"
)

func getter(name string) (string, int) {
	return name, len(name)
}

type My struct {
	Name   string
	Power  int
	Father *My
}

func NewMy(name string, power int) My {
	return My{
		Name:  name,
		Power: power,
	}
}

func Sup(m *My) {
	m.Power *= 10
}

func (s *My) Super() {
	s.Power += 10
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Недостаточно аргументов")
		os.Exit(1)
	}
	// Объявили экземпляр
	my1 := new(My)
	my1.Name, my1.Power = getter(os.Args[1])
	// Создали экземпляр
	my2 := NewMy(getter(os.Args[1]))
	Sup(my1)
	my2.Super()
	// Присвоили ссылку на отца
	my1.Father = &my2
	fmt.Println("Тут вывод: \n", my1.Father, my2)
}
