package main

import (
	"hangman"
)

func main() {
	files := []string{"words.txt", "words2.txt", "words3.txt"}
	game := hangman.NewGame(files)
	game.Play()
}
