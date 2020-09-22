package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

type deck []string

func newDeck() (cards deck) {
	cardSuits := []string{"Spades", "Diamonds", "Hearts", "Clubs"}
	cardValues := []string{"Ace", "Two", "Three", "Four", "Five",
		"Six", "Seven", "Eight", "Nine", "Ten", "Queen", "Jack", "King"}

	for _, value := range cardValues {
		for _, suit := range cardSuits {
			cards = append(cards, value+" of "+suit)
		}
	}
	return
}

func (d deck) print() {
	for i, card := range d {
		fmt.Println(i+1, card)
	}
}

func draw(d deck, handSize int) (deck, deck) {
	return d[:handSize], d[handSize:]
}

func (d deck) toString() string {
	return strings.Join([]string(d), ",")
}

func (d deck) shuffle() {
	source := rand.NewSource(int64(time.Now().UnixNano()))
	rng := rand.New(source)

	for i := range d {
		newPos := rng.Intn(len(d) - 1)
		d[i], d[newPos] = d[newPos], d[i]
	}
}

func (d deck) saveToFile(filename string) error {
	return ioutil.WriteFile(filename, []byte(d.toString()), 0666)
}

func newDeckFromFile(filename string) deck {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Failed to read deck from file |", err)
		fmt.Println("Creating new deck...")
		return newDeck()
	}

	fmt.Println("Successfully read a deck from file!")
	s := strings.Split(string(bs), ",")
	return deck(s)
}
