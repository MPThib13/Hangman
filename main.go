package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func clearTerminal() {
	fmt.Print("\033[H\033[2J")
}

func drawHangman(errors int) {
	stages := []string{
		`
		______
		|    |
		|    
		|   
		|    
		|___
		`,
		`
		______
		|    |
		|    O
		|   
		|    
		|___
		`,
		`
		______
		|    |
		|    O
		|    |
		|    
		|___
		`,
		`
		______
		|    |
		|    O
		|   /|
		|    
		|___
		`,
		`
		______
		|    |
		|    O
		|   /|\
		|    
		|___
		`,
		`
		______
		|    |
		|    O
		|   /|\
		|   / 
		|___
		`,
		`
		______
		|    |
		|    O
		|   /|\
		|   / \
		|___
		`,
	}
	fmt.Println(stages[errors])
}

func main() {
	file, err := os.Open("hangman.txt")
	if err != nil {
		log.Fatalf("Impossible d'ouvrir le fichier : %v", err)
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		words = append(words, strings.Fields(line)...)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Erreur lors de la lecture du fichier : %v", err)
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(words))
	selectedWord := words[randomIndex]
	wordLength := len(selectedWord)

	hiddenWord := make([]rune, wordLength)
	for i := 0; i < wordLength; i++ {
		hiddenWord[i] = '_'
	}

	if wordLength > 10 {
		firstLetterIndex := rand.Intn(wordLength)
		secondLetterIndex := rand.Intn(wordLength)
		for firstLetterIndex == secondLetterIndex {
			secondLetterIndex = rand.Intn(wordLength)
		}
		hiddenWord[firstLetterIndex] = rune(selectedWord[firstLetterIndex])
		hiddenWord[secondLetterIndex] = rune(selectedWord[secondLetterIndex])
	} else {
		letterIndex := rand.Intn(wordLength)
		hiddenWord[letterIndex] = rune(selectedWord[letterIndex])
	}

	maxErrors := 6
	currentErrors := 0

	clearTerminal()
	drawHangman(currentErrors)
	fmt.Printf("Mot à deviner : %s\n", string(hiddenWord))

	reader := bufio.NewReader(os.Stdin)
	guessed := false

	for !guessed {
		fmt.Print("Entrez une lettre ou le mot entier : ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		input = strings.ToLower(input)

		actionPerformed := false

		if len(input) == len(selectedWord) {
			if input == strings.ToLower(selectedWord) {
				guessed = true
				clearTerminal()
				drawHangman(currentErrors)
				fmt.Printf("Félicitations ! Vous avez trouvé le mot : %s\n", selectedWord)
				break
			} else {
				fmt.Println("Ce n'est pas le bon mot.")
				currentErrors++
				actionPerformed = true
			}
		} else if len(input) == 1 {
			letter := rune(input[0])
			found := false

			for i, char := range selectedWord {
				if char == letter {
					hiddenWord[i] = letter
					found = true
				}
			}

			clearTerminal()

			if found {
				fmt.Println("Bonne lettre !")
			} else {
				fmt.Println("Cette lettre n'est pas dans le mot.")
				currentErrors++
			}
			actionPerformed = true

			if actionPerformed {
				drawHangman(currentErrors)
				fmt.Printf("Mot à deviner : %s\n", string(hiddenWord))
			}
		} else {
			fmt.Println("Veuillez entrer soit une lettre, soit le mot entier.")
		}

		guessed = true
		for _, char := range hiddenWord {
			if char == '_' {
				guessed = false
				break
			}
		}

		if guessed {
			clearTerminal()
			drawHangman(currentErrors)
			fmt.Printf("Félicitations ! Vous avez deviné le mot : %s\n", selectedWord)
		}

		if currentErrors >= maxErrors {
			clearTerminal()
			drawHangman(currentErrors)
			fmt.Printf("Vous avez perdu ! Le mot était : %s\n", selectedWord)
			break
		} else {
			fmt.Printf("Erreurs restantes : %d\n", maxErrors-currentErrors)
		}
	}
}
