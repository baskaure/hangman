package hangman

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type Game struct {
	Words           []string
	MotAleatoire    string
	LettresRevelees map[int]bool
	Tentatives      int
	Positions       []int
	MotAffiche      string
}

func (g *Game) Display() {
	motAffiche := ""
	for i, lettre := range g.MotAleatoire {
		if g.LettresRevelees[i] {
			motAffiche += strings.ToUpper(string(lettre)) + " "
		} else {
			motAffiche += "_ "
		}
	}
	g.MotAffiche = motAffiche
}

func NewGame() *Game {
	rand.Seed(time.Now().Unix())
	var words []string
	var fileName string

	fileName = "words.txt"
	file, err := os.Open(fileName)
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	motAleatoire := strings.ToUpper(words[rand.Intn(len(words))])
	n := len(motAleatoire)/2 - 1
	lettresRevelees := make(map[int]bool)
	for i := 0; i < n; i++ {
		randIndex := rand.Intn(len(motAleatoire))
		lettresRevelees[randIndex] = true
	}
	fmt.Println("test")
	return &Game{
		Words:           words,
		MotAleatoire:    motAleatoire,
		LettresRevelees: lettresRevelees,
		Tentatives:      10,
	}
}

func (g *Game) Play(w http.ResponseWriter, r *http.Request, choice string) {
	g.Positions = []int{72, 64, 56, 48, 40, 32, 24, 16, 8, 0}

	var lettresSuggerees []string

	if choice == "STOP" || choice == "stop" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if len(choice) == 1 && IsLetter(choice) {
		lettre := strings.ToUpper(choice)

		if Contains(lettresSuggerees, lettre) {
			fmt.Printf("You already suggested the letter '%s'. Try again.\n", lettre)
			return
		}

		lettresSuggerees = append(lettresSuggerees, lettre)

		lettreTrouvée := false
		for i, lettreMot := range g.MotAleatoire {
			if lettre == string(lettreMot) && !g.LettresRevelees[i] {
				g.LettresRevelees[i] = true
				lettreTrouvée = true
			}
		}

		if !lettreTrouvée {
			g.Tentatives--
			fmt.Printf("Not present in the word, %d attempts remaining\n", g.Tentatives)
		}
	} else if len(choice) >= 2 {
		mot := strings.ToUpper(choice)

		if mot == g.MotAleatoire {
			g.LettresRevelees = make(map[int]bool)
			for i := range g.MotAleatoire {
				g.LettresRevelees[i] = true
			}
			g.Display()
			return
		} else {
			g.Tentatives -= 2
			fmt.Printf("Incorrect word, %d attempts remaining\n", g.Tentatives)
		}
	} else {
		fmt.Println("Invalid input. Please enter a single letter or a word of at least two characters.")
	}

	motDevine := true
	for i := range g.MotAleatoire {
		if !g.LettresRevelees[i] {
			motDevine = false
			break
		}
	}
	if motDevine {
		g.Display()
		return
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
