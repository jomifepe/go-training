package main

import "fmt"

type person struct {
	firstName string
	lastName  string
	contact
}

type contact struct {
	email   string
	zipCode int
}

func main() {
	eliot := person{
		firstName: "Elliot",
		lastName:  "Alderson",
		contact: contact{
			email:   "eliot@gmail.com",
			zipCode: 94000,
		},
	}
	eliot.print()
	eliot.updateName("John")
	eliot.print()
}

func (p person) print() {
	fmt.Printf("This is %v %v, and it's contact is %v, %v\n", p.firstName, p.lastName, p.email, p.zipCode)
}

func (p *person) updateName(newFirstName string) {
	(*p).firstName = newFirstName
}
