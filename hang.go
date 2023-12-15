package hangman

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Game struct {
	words           []string
	motAleatoire    string
	lettresRevelees map[int]bool
	tentatives      int
	positions       []int
	ascii           string
}

func (g *Game) Restore() error {
	file, err := os.Open("save.txt")
	if err != nil {
		return fmt.Errorf("Error opening save file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(g); err != nil {
		return fmt.Errorf("Error decoding save file: %v", err)
	}

	fmt.Println("Game restored successfully.")
	return nil
}

func NewGame() *Game {
	rand.Seed(time.Now().Unix())

	var words []string
	var fileName string

	println()
	fileName = "../words.txt"

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

	return &Game{
		words:           words,
		motAleatoire:    motAleatoire,
		lettresRevelees: lettresRevelees,
		tentatives:      10,
	}
}

func (g *Game) Save() {
	save, err := json.Marshal(g)
	if err != nil {
		fmt.Println("Error marshaling game:", err)
		os.Exit(1)
	}

	file, err := os.Create("save.txt")
	if err != nil {
		fmt.Println("Error creating save file:", err)
		os.Exit(2)
	}
	defer file.Close()

	_, err = file.Write(save)
	if err != nil {
		fmt.Println("Error writing to save file:", err)
		os.Exit(3)
	}
}

func (g *Game) Display() {
	reset := "\033[0m"
	gray := "\033[30m"

	motAffiche := ""
	for i, lettre := range g.motAleatoire {
		if g.lettresRevelees[i] {
			motAffiche += strings.ToUpper(string(lettre)) + " "
		} else {
			motAffiche += "_ "
		}
	}
	fmt.Println(gray + "---------------------------------" + reset)
	fmt.Printf("%s\n", motAffiche)
	fmt.Println(gray + "---------------------------------" + reset)
}

func (g *Game) Play() {
	g.positions = []int{72, 64, 56, 48, 40, 32, 24, 16, 8, 0}
	green := "\033[32m"
	red := "\033[31m"
	gray := "\033[30m"
	reset := "\033[0m"

	var lettresSuggerees []string

	for g.tentatives > 0 {
		g.Display()
		fmt.Print("Choose: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := strings.ToUpper(scanner.Text())

		if input == "STOP" {
			fmt.Println("To resume the game you must --startWith save.txt")
			g.Save()
			return
		}

		if len(input) == 1 && isLetter(input) {
			lettre := strings.ToUpper(input)

			if contains(lettresSuggerees, lettre) {
				fmt.Printf("You already suggested the letter '%s'. Try again.\n", lettre)
				continue
			}

			lettresSuggerees = append(lettresSuggerees, lettre)

			lettreTrouvée := false
			for i, lettreMot := range g.motAleatoire {
				if lettre == string(lettreMot) && !g.lettresRevelees[i] {
					g.lettresRevelees[i] = true
					lettreTrouvée = true
				}
			}

			if !lettreTrouvée {
				g.tentatives--
				fmt.Printf("Not present in the word, %d attempts remaining\n", g.tentatives)
				Positions := Positions(g)
				fmt.Println(Positions)
			}
		} else if len(input) >= 2 {
			mot := strings.ToUpper(input)

			if mot == g.motAleatoire {
				g.lettresRevelees = make(map[int]bool)
				for i := range g.motAleatoire {
					g.lettresRevelees[i] = true
				}
				g.Display()
				fmt.Printf(green + "Congrats !\n")
				fmt.Println(gray + "---------------------------------" + reset)
				fmt.Println("          END OF GAME")
				fmt.Println(gray + "---------------------------------" + reset)
				return
			} else {
				g.tentatives -= 2
				fmt.Printf("Incorrect word, %d attempts remaining\n", g.tentatives)
				Positions := Positions(g)
				fmt.Println(Positions)
			}
		} else {
			fmt.Println("Invalid input. Please enter a single letter or a word of at least two characters.")
		}

		motDevine := true
		for i := range g.motAleatoire {
			if !g.lettresRevelees[i] {
				motDevine = false
				break
			}
		}
		if motDevine {
			g.Display() // Afficher le mot à deviner
			fmt.Printf(green + "Congrats !\n")
			fmt.Println(gray + "---------------------------------" + reset)
			fmt.Println("          END OF GAME")
			fmt.Println(gray + "---------------------------------" + reset)
			return // Quitter la fonction
		}
	}
	fmt.Println(gray + "---------------------------------" + reset)
	fmt.Printf(red+"Sorry, you didn't guess the word. The word was '%s'.\n", g.motAleatoire)
	fmt.Println(gray + "---------------------------------" + reset)
	fmt.Println("          END OF GAME")
	fmt.Println(gray + "---------------------------------" + reset)
}

func isLetter(s string) bool {
	return len(s) == 1 && ('a' <= s[0] && s[0] <= 'z' || 'A' <= s[0] && s[0] <= 'Z')
}

func contains(liste []string, lettre string) bool {
	for _, l := range liste {
		if l == lettre {
			return true
		}
	}
	return false
}

func Positions(g *Game) string {
	hangmandraw := ""
	file, err := os.Open("hangman.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	line := 0
	if g.tentatives < len(g.positions) {
		for scanner.Scan() {
			line++
			if line > g.positions[g.tentatives] && line <= g.positions[g.tentatives]+7 {
				hangmandraw += "\n" + scanner.Text()
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return hangmandraw
}
