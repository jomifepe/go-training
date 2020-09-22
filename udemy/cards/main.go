package main

const filename string = "my_cards"

func main() {
	cards := newDeck()
	cards.shuffle()
	cards.print()
	cards.saveToFile(filename)
}
