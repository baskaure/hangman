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
	ResetPart        bool
}

func Reset(g *Game) {
	g.ResetPart = true
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

func NewGame(g *Game) {
	g.ResetPart = false
	rand.Seed(time.Now().Unix())
	words := LoadDictionary("words.txt")
	motAleatoire := strings.ToUpper(words[rand.Intn(len(words))])
	g.Words = words
	g.MotAleatoire = motAleatoire
	g.LettresRevelees = make(map[int]bool)

	initialRevealedLetters := 2
	for i := 0; i < initialRevealedLetters; i++ {
		randIndex := rand.Intn(len(motAleatoire))
		g.LettresRevelees[randIndex] = true
	}

	g.Tentatives = 10
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

	allLettersFound := true
	for i, letterMot := range g.MotAleatoire {
		if letter == string(letterMot) && !g.LettresRevelees[i] {
			g.LettresRevelees[i] = true
		}
		if !g.LettresRevelees[i] {
			allLettersFound = false
		}
	}

	if allLettersFound {
		g.FoundWord = 1
	}

	if !allLettersFound {
		g.Tentatives--
		g.Message = "Pas présent dans le mot !"
	}

	Display(g)

	if g.Tentatives <= 0 || g.FoundWord == 1 {
		if g.FoundWord == 1 {
			g.Message = "Félicitations ! Vous avez trouvé le mot: " + g.MotAleatoire
		} else {
			g.Message = "Désolé, vous avez épuisé toutes vos tentatives. Le mot était: " + g.MotAleatoire
		}
		Reset(g)
	}
}

func PlayWord(g *Game, word string) {
	if word == g.MotAleatoire {
		for i := range g.MotAleatoire {
			g.LettresRevelees[i] = true
		}
		g.FoundWord = 1
		Display(g)
		return
	}

	g.Tentatives -= 2
	g.Message = "Mot incorrect !"

	if g.Tentatives <= 0 || g.FoundWord == 1 {
		if g.FoundWord == 1 {
			g.Message = "Félicitations ! Vous avez trouvé le mot: " + g.MotAleatoire
		} else {
			g.Message = "Désolé, vous avez épuisé toutes vos tentatives. Le mot était: " + g.MotAleatoire
		}
		Reset(g)
	}
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
