package hangman

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Game struct {
	Words            []string
	RandomWord       string
	RevealedLetters  map[int]bool
	Tries            int
	Positions        []int
	DisplayedWord    string
	SuggestedLetters []string
	FoundWord        int
	Message          string
	MessageReveal    string
}

func Reset(g *Game) {
	g.RevealedLetters = make(map[int]bool)
	initialRevealedLetters := 2
	for i := 0; i < initialRevealedLetters; i++ {
		randIndex := rand.Intn(len(g.RandomWord))
		g.RevealedLetters[randIndex] = true
	}

	g.Tries = 10
	g.SuggestedLetters = nil
	g.FoundWord = 0
	g.Message = ""
	g.MessageReveal = ""
}

const (
	EasyModeFile   = "words.txt"
	NormalModeFile = "words2.txt"
	HardModeFile   = "words3.txt"
)

func NewGame(g *Game, mode string) {
	rand.Seed(time.Now().Unix())

	var words []string
	var randomWord string

	switch mode {
	case "easy":
		words = LoadDictionary(EasyModeFile)
	case "normal":
		words = LoadDictionary(NormalModeFile)
	case "hard":
		words = LoadDictionary(HardModeFile)
	default:
		words = LoadDictionary(EasyModeFile)
	}

	randomWord = strings.ToUpper(words[rand.Intn(len(words))])

	g.Words = words
	g.RandomWord = randomWord
	g.RevealedLetters = make(map[int]bool)

	initialRevealedLetters := 2
	for i := 0; i < initialRevealedLetters; i++ {
		randIndex := rand.Intn(len(randomWord))
		g.RevealedLetters[randIndex] = true
	}

	g.Tries = 10
	g.FoundWord = 0
}

func Display(g *Game) {
	displayedWord := ""
	for i, letter := range g.RandomWord {
		if g.RevealedLetters[i] {
			displayedWord += strings.ToUpper(string(letter)) + " "
		} else {
			displayedWord += "_ "
		}
	}
	g.DisplayedWord = displayedWord
	g.MessageReveal = strings.Join(g.SuggestedLetters, ", ")
}

func Play(g *Game, choice string) {
	if len(choice) == 1 && IsLetter(choice) {
		PlayLetter(g, strings.ToUpper(choice))
	} else if len(choice) >= 2 {
		PlayWord(g, strings.ToUpper(choice))
	} else {
		g.Message = "Invalid input. Please enter a single letter or a word of at least two characters."
	}
}

func PlayLetter(g *Game, letter string) {
	g.FoundWord = 0

	if Contains(g.SuggestedLetters, letter) {
		g.Message = "You have already suggested the letter."
		return
	}

	g.SuggestedLetters = append(g.SuggestedLetters, letter)

	letterFound := false
	for i, letterWord := range g.RandomWord {
		if letter == string(letterWord) && !g.RevealedLetters[i] {
			g.RevealedLetters[i] = true
			letterFound = true
		}
	}

	if !letterFound {
		g.Tries--
		g.Message = "Not present in the word!"
	} else {
		wordFoundCount := 0
		for _, revealed := range g.RevealedLetters {
			if revealed {
				wordFoundCount++
			}
		}

		if wordFoundCount == len(g.RandomWord) {
			g.FoundWord = 1
		}
	}

	Display(g)
}

func PlayWord(g *Game, word string) {
	if g.FoundWord == 1 || g.Tries <= 0 {
		Reset(g)
		return
	}

	if word == g.RandomWord {
		for i := range g.RandomWord {
			g.RevealedLetters[i] = true
		}
		g.Message = "Congratulations, you found the word!"
		Display(g)
		g.FoundWord = 1
		return
	}

	g.Tries -= 2
	g.Message = "Incorrect word!"

	wordFoundCount := 0
	for _, revealed := range g.RevealedLetters {
		if revealed {
			wordFoundCount++
		}
	}

	if wordFoundCount == len(g.RandomWord) {
		g.FoundWord = 1
	}

	Display(g)
}

func IsLetter(s string) bool {
	return len(s) == 1 && ('a' <= s[0] && s[0] <= 'z' || 'A' <= s[0] && s[0] <= 'Z')
}

func Contains(list []string, letter string) bool {
	for _, l := range list {
		if l == letter {
			return true
		}
	}
	return false
}

func LoadDictionary(fileName string) []string {
	file, err := os.Open(fileName)
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	return words
}
