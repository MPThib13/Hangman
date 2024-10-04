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

func loadHangmanDrawings(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var drawings []string
	scanner := bufio.NewScanner(file)
	var currentDrawing strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if currentDrawing.Len() > 0 {
				drawings = append(drawings, currentDrawing.String())
				currentDrawing.Reset()
			}
		} else {
			currentDrawing.WriteString(line + "\n")
		}
	}

	if currentDrawing.Len() > 0 {
		drawings = append(drawings, currentDrawing.String())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return drawings, nil
}

func drawHangman(errors int, drawings []string) {
	if errors < len(drawings) {
		fmt.Println(drawings[errors])
	}
}

func listTextFiles() ([]string, error) {
	dir, err := os.Open(".")
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	var textFiles []string
	for _, file := range files {
		if strings.HasSuffix(file, ".txt") && file != "hangman-drawing.txt" {
			textFiles = append(textFiles, file)
		}
	}

	return textFiles, nil
}

func chooseFile(files []string) string {
	if len(files) == 1 {
		fmt.Printf("Un seul fichier trouvé : %s. Utilisation de ce fichier.\n", files[0])
		return files[0]
	}

	fmt.Println("Fichiers disponibles :")
	for index, file := range files {
		fmt.Printf("%d. %s\n", index+1, file)
	}

	var choice int
	for {
		fmt.Print("Choisissez le numéro du fichier à utiliser : ")
		_, err := fmt.Scanf("%d", &choice)
		if err == nil && choice > 0 && choice <= len(files) {
			break
		}
		fmt.Println("Numéro invalide, veuillez réessayer.")
	}

	return files[choice-1]
}

func main() {
	hangmanDrawings, err := loadHangmanDrawings("hangman-drawing.txt")
	if err != nil {
		log.Fatalf("Erreur lors du chargement des dessins du pendu : %v", err)
	}

	files, err := listTextFiles()
	if err != nil || len(files) == 0 {
		log.Fatalf("Aucun fichier .txt trouvé dans le répertoire courant ou erreur lors de la lecture : %v", err)
	}

	selectedFile := chooseFile(files)

	file, err := os.Open(selectedFile)
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
	usedLetters := make(map[rune]bool)

	clearTerminal()
	drawHangman(currentErrors, hangmanDrawings)
	fmt.Printf("Mot à deviner : %s\n", string(hiddenWord))
	fmt.Println("Lettres déjà utilisées : Aucun")

	reader := bufio.NewReader(os.Stdin)
	guessed := false

	for !guessed {
		fmt.Print("Entrez une lettre ou le mot entier : ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		input = strings.ToLower(input)

		if len(input) == 1 {
			letter := rune(input[0])
			if usedLetters[letter] {
				fmt.Printf("Vous avez déjà utilisé la lettre '%c'. Essayez une autre.\n", letter)
				continue
			} else {
				usedLetters[letter] = true
			}
		}

		actionPerformed := false

		if len(input) == len(selectedWord) {
			if input == strings.ToLower(selectedWord) {
				guessed = true
				clearTerminal()
				drawHangman(currentErrors, hangmanDrawings)
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
				drawHangman(currentErrors, hangmanDrawings)
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
			drawHangman(currentErrors, hangmanDrawings)
			fmt.Printf("Félicitations ! Vous avez deviné le mot : %s\n", selectedWord)
		}

		if currentErrors >= maxErrors {
			clearTerminal()
			drawHangman(currentErrors, hangmanDrawings)
			fmt.Printf("Vous avez perdu ! Le mot était : %s\n", selectedWord)
			break
		} else {
			fmt.Printf("Erreurs restantes : %d\n", maxErrors-currentErrors)
			fmt.Printf("Lettres déjà utilisées : ")
			for letter := range usedLetters {
				fmt.Printf("%c ", letter)
			}
			fmt.Println()
		}
	}
}
