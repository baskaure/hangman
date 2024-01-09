package hangman

import (
	"bufio"
	"fmt"
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
}

func NewGame(g *Game) {
	rand.Seed(time.Now().Unix())
	words := LoadDictionary("words.txt")
	motAleatoire := strings.ToUpper(words[rand.Intn(len(words))])
	lettresRevelees := make(map[int]bool)
	n := len(motAleatoire)/2 - 1
	for i := 0; i < n; i++ {
		randIndex := rand.Intn(len(motAleatoire))
		lettresRevelees[randIndex] = true
	}
	g.Words = words
	g.MotAleatoire = motAleatoire
	g.LettresRevelees = lettresRevelees
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
	fmt.Printf("Mot à deviner: %s\n", g.MotAffiche)
	fmt.Printf("Lettres utilisées: %s\n", strings.Join(g.LettresSuggerees, ", "))
	fmt.Printf("Tentatives restantes: %d\n", g.Tentatives)
}

func Play(g *Game, choice string) {
	if len(choice) == 1 && IsLetter(choice) {
		PlayLetter(g, strings.ToUpper(choice))
	} else if len(choice) >= 2 {
		PlayWord(g, strings.ToUpper(choice))
	} else {
		fmt.Println("Entrée non valide. Veuillez entrer une lettre unique ou un mot d'au moins deux caractères.")
	}
}

func PlayLetter(g *Game, letter string) {
	if Contains(g.LettresSuggerees, letter) {
		fmt.Printf("Vous avez déjà proposé la lettre '%s'. Réessayez.\n", letter)
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
		fmt.Printf("Pas présent dans le mot, %d tentatives restantes\n", g.Tentatives)
	}
}

func PlayWord(g *Game, word string) {
	if word == g.MotAleatoire {
		for i := range g.MotAleatoire {
			g.LettresRevelees[i] = true
		}
		Display(g)
		return
	}

	g.Tentatives -= 2
	fmt.Printf("Mot incorrect, %d tentatives restantes\n", g.Tentatives)
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
