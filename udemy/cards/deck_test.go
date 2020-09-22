package main

import (
	"os"
	"strings"
	"testing"
)

func TestNewDeck(t *testing.T) {
	d, expectedSize := newDeck(), 52
	if len(d) != expectedSize {
		t.Errorf("Expected deck length of %v, but got %v", expectedSize, len(d))
	}
	if strings.ToLower(d[0]) != "ace of spades" {
		t.Errorf("Expected first card of ace of spades, but got %v", d[0])
	}

	if strings.ToLower(d[len(d)-1]) != "king of clubs" {
		t.Errorf("Expected last card of king of clubs, but got %v", d[len(d)-1])
	}
}

func TestSaveToDeckAndNewDeckFromFile(t *testing.T) {
	filename := "_decktesting"
	os.Remove(filename)
	deck := newDeck()
	deck.saveToFile(filename)

	loadedDeck := newDeckFromFile(filename)
	if len(loadedDeck) != 52 {
		t.Errorf("Expected 52 cards in deck, got %v", len(loadedDeck))
	}

	os.Remove(filename)
}
