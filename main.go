package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pannipasra/pokedexcli/internals/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(config *pokeapi.Config, client *pokeapi.Client) error
}

var commandLists map[string]cliCommand

func main() {
	// Create a scanner that reads from standard input (os.Stdin)
	scanner := bufio.NewScanner(os.Stdin)
	prompt := "Pokedex > "

	// Initiate PokeAPI client and config
	client := pokeapi.NewClient()
	config := &pokeapi.Config{
		Next:     nil,
		Previous: nil,
	}

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
			description: "Displays 20 Pokémon next locations per map call.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays 20 Pokémon previous locations per map call.",
			callback:    commandMapb,
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
					err := command.callback(config, client)
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

func commandExit(config *pokeapi.Config, client *pokeapi.Client) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil // This line will never execute due to os.Exit
}

func commandHelp(config *pokeapi.Config, client *pokeapi.Client) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	for key, value := range commandLists {
		fmt.Printf("%v: %v\n", key, value.description)
	}

	return nil
}

func commandMap(config *pokeapi.Config, client *pokeapi.Client) error {
	res, err := client.ListLocationAreas(config)
	if err != nil {
		return err
	}

	// Print the results
	for _, result := range res.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandMapb(config *pokeapi.Config, client *pokeapi.Client) error {
	res, err := client.ListPreviousLocationAreas(config)
	if err != nil {
		return err
	}
	// Print the results
	for _, result := range res.Results {
		fmt.Println(result.Name)
	}

	return nil
}
