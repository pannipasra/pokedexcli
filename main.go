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
	callback    func(client *pokeapi.Client, config *pokeapi.Config, param string) error
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
		"explore": {
			name:        "explore",
			description: "Explore a location area for Pokémon. Usage: explore <area_name>",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catching Pokemon by name. Usage: catch <pokemon_name>",
			callback:    commandCatch,
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

			commandName := inputs[0]
			param := ""
			if len(inputs) > 1 {
				param = inputs[1]
			}

			// Check if the first word is a command
			if command, exists := commandLists[commandName]; exists {
				// Command exists, execute its callback
				err := command.callback(client, config, param)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error executing command:", err)
				}
			} else {
				fmt.Println("Unknown command")
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

func commandExit(client *pokeapi.Client, config *pokeapi.Config, param string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil // This line will never execute due to os.Exit
}

func commandHelp(client *pokeapi.Client, config *pokeapi.Config, param string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	for key, value := range commandLists {
		fmt.Printf("%v: %v\n", key, value.description)
	}

	return nil
}

func commandMap(client *pokeapi.Client, config *pokeapi.Config, param string) error {
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

func commandMapb(client *pokeapi.Client, config *pokeapi.Config, param string) error {
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

func commandExplore(client *pokeapi.Client, config *pokeapi.Config, locationName string) error {
	if locationName == "" {
		return fmt.Errorf("area name is required. Usage: explore <area_name>")
	}

	exploreEncounter, err := client.Explore(locationName)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", locationName)
	fmt.Println("Found Pokemon:")
	for _, pokemon := range exploreEncounter.PokemonEncounters {
		fmt.Println("-", pokemon.Pokemon.Name)
	}

	return nil
}

func commandCatch(client *pokeapi.Client, config *pokeapi.Config, pokemonName string) error {
	if pokemonName == "" {
		return fmt.Errorf("pokemon name is required. Usage: catch <pokemon_name>")
	}

	pokemon, err := client.Catch(pokemonName)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	fmt.Printf("%s escaped!\n", pokemonName)
	fmt.Printf("%s was caught!\n", pokemonName)
	fmt.Printf("%s has base_experience %v\n", pokemonName, pokemon.BaseExperience)

	return nil
}
