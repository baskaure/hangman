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
	MotAleatoire     string
	LettresRevelees  map[int]bool
	Tentatives       int
	Positions        []int
	MotAffiche       string
	LettresSuggerees []string
	FoundWord        int
	Message          string
	MessageReveal    string
}

func Reset(g *Game) {
	g.LettresRevelees = make(map[int]bool)
	initialRevealedLetters := 2
	for i := 0; i < initialRevealedLetters; i++ {
		randIndex := rand.Intn(len(g.MotAleatoire))
		g.LettresRevelees[randIndex] = true
	}

	g.Tentatives = 10
	g.LettresSuggerees = nil
	g.FoundWord = 0
	g.Message = ""
	g.MessageReveal = ""
	NewGame(g)
	Display(g)
}

const (
	EasyModeFile   = "words.txt"
	NormalModeFile = "words2.txt"
	HardModeFile   = "words3.txt"
)

func NewGame(g *Game, mode string) {
	rand.Seed(time.Now().Unix())

	var words []string
	var motAleatoire string

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

	motAleatoire = strings.ToUpper(words[rand.Intn(len(words))])

	g.Words = words
	g.MotAleatoire = motAleatoire
	g.LettresRevelees = make(map[int]bool)

	initialRevealedLetters := 2
	for i := 0; i < initialRevealedLetters; i++ {
		randIndex := rand.Intn(len(motAleatoire))
		g.LettresRevelees[randIndex] = true
	}

	g.Tentatives = 10
	g.FoundWord = 0
}

func Display(g *Game) {
	motAffiche := ""
	for i, lettre := range g.MotAleatoire {
		if g.LettresRevelees[i] {
			motAffiche += strings.ToUpper(string(lettre)) + " "
		} else {
			motAffiche += "_ "
		}
	}
	g.MotAffiche = motAffiche
	g.MessageReveal = strings.Join(g.LettresSuggerees, ", ")
}

func Play(g *Game, choice string) {
	if len(choice) == 1 && IsLetter(choice) {
		PlayLetter(g, strings.ToUpper(choice))
	} else if len(choice) >= 2 {
		PlayWord(g, strings.ToUpper(choice))
	} else {
		g.Message = "Entrée non valide. Veuillez entrer une lettre unique ou un mot d'au moins deux caractères."
	}
}

func PlayLetter(g *Game, letter string) {
	g.FoundWord = 0

	if Contains(g.LettresSuggerees, letter) {
		g.Message = "Vous avez déjà proposé la lettre."
		return
	}

	g.LettresSuggerees = append(g.LettresSuggerees, letter)

	letterFound := false
	for i, letterMot := range g.MotAleatoire {
		if letter == string(letterMot) && !g.LettresRevelees[i] {
			g.LettresRevelees[i] = true
			letterFound = true
		}
	}

	if !letterFound {
		g.Tentatives--
		g.Message = "Pas présent dans le mot !"
	} else {
		wordFoundCount := 0
		for _, revealed := range g.LettresRevelees {
			if revealed {
				wordFoundCount++
			}
		}

		if wordFoundCount == len(g.MotAleatoire) {
			g.FoundWord = 1
		}
	}

	Display(g)
}

func PlayWord(g *Game, word string) {
	if g.FoundWord == 1 || g.Tentatives <= 0 {
		Reset(g)
		return
	}

	if word == g.MotAleatoire {
		for i := range g.MotAleatoire {
			g.LettresRevelees[i] = true
		}
		g.Message = "Félicitations, vous avez trouvé le mot!"
		Display(g)
		g.FoundWord = 1
		return
	}

	g.Tentatives -= 2
	g.Message = "Mot incorrect !"

	wordFoundCount := 0
	for _, revealed := range g.LettresRevelees {
		if revealed {
			wordFoundCount++
		}
	}

	if wordFoundCount == len(g.MotAleatoire) {
		g.FoundWord = 1
	}

	Display(g)
}

func IsLetter(s string) bool {
	return len(s) == 1 && ('a' <= s[0] && s[0] <= 'z' || 'A' <= s[0] && s[0] <= 'Z')
}

func Contains(liste []string, lettre string) bool {
	for _, l := range liste {
		if l == lettre {
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
