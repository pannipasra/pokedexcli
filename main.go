package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
)

func main(){
	prompt := "Pokedex > "

	for {
	// Create a scanner that reads from standard input (os.Stdin)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print(prompt)

	// Use Scan() to read the next line of input
	if scanner.Scan() {
		// Get text that was read
		input := scanner.Text()

		cleanInputWords := cleanInput(input)

		fmt.Println("Your command was: ", cleanInputWords[0])
	}

	// Check for erros
	if err := scanner.Err(); err != nil {
		fmt.Println(os.Stderr, "Error reading input:", err)
	}


}
}

func cleanInput(text string) []string {
	// Trim leading and trailing whitespace
	text = strings.TrimSpace(text)
	
	text = strings.ToLower(text)

	// Convert to lowercase and split by whitespace
	words := strings.Fields(text)


	return words
}

