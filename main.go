package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type PokedexMap struct {
	Count    int                `json:"count"`
	Next     string             `json:"next"`
	Previous interface{}        `json:"previous"` // Using interface{} to handle null
	Results  []PokedexMapResult `json:"results"`
}

type PokedexMapResult struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

var commandLists map[string]cliCommand

func main() {
	// Create a scanner that reads from standard input (os.Stdin)
	scanner := bufio.NewScanner(os.Stdin)
	prompt := "Pokedex > "
	commandLists = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays 20 Pokémon locations per map call.",
			callback:    commandMap,
		},
	}

	for {
		fmt.Print(prompt)

		// Use Scan() to read the next line of input
		if scanner.Scan() {
			// Get text that was read
			input := scanner.Text()

			inputs := cleanInput(input)

			if len(inputs) == 0 {
				continue
			}

			for _, commandName := range inputs {
				// Check if the first word is a command
				if command, exists := commandLists[commandName]; exists {
					// Command exists, execute its callback
					err := command.callback()
					if err != nil {
						fmt.Fprintln(os.Stderr, "Error executing command:", err)
					}
				} else {
					fmt.Println("Unknown command")
				}
			}
		}

		// Check for errors
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
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

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil // This line will never execute due to os.Exit
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	for key, value := range commandLists {
		fmt.Printf("%v: %v\n", key, value.description)
	}

	return nil
}

func commandMap() error {
	url := "https://pokeapi.co/api/v2/location-area?offset=20&limit=20"

	// Make http request
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Read response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// Parse JSON into PokedexMap struct
	var pokedex PokedexMap
	err = json.Unmarshal(body, &pokedex)
	if err != nil {
		return err
	}

	fmt.Printf("Count: %d\n", pokedex.Count)
	fmt.Printf("Next: %s\n", pokedex.Next)
	fmt.Printf("Previous: %v\n", pokedex.Previous)

	for i, result := range pokedex.Results {
		fmt.Printf("Result %d: %s - %s\n", i+1, result.Name, result.Url)
	}

	return nil
}
