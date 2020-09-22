package main

import "fmt"

type square struct {
	sideLength float64
}
type triangle struct {
	height float64
	base   float64
}

type shape interface {
	getArea() float64
}

func main() {
	fmt.Print("Square area: ")
	printArea(square{4})
	fmt.Print("Triangle area: ")
	printArea(triangle{height: 4, base: 3})
}

func printArea(s shape) {
	fmt.Println(s.getArea())
}

func (s square) getArea() float64 {
	return s.sideLength * s.sideLength
}

func (t triangle) getArea() float64 {
	return .5 * t.base * t.height
}
