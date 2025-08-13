// main.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/Marcus-Gustafsson/pokedexCLI/internal"
)

// main starts the Pokedex REPL (Read-Eval-Print Loop).
// It displays a prompt, reads user commands, and dispatches them to the proper handler.
// The loop continues until standard input ends or the user issues an exit command.
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Pokedex > ")

	// configPTR keeps track of paging state for the PokeAPI.
	configPTR := config{}
	pokedex := make(map[string]pokemonDetails)

	// Init new cache with given interval (interval determines when cacheEntries are cleared)
	cachePtr := internal.NewCache(30 * time.Second)

	// The REPL loop: waits for user input, dispatches commands, then re-prompts.
	for scanner.Scan() {
		userInput := scanner.Text()
		cleanedWords := cleanInput(userInput)

		if len(cleanedWords) > 0 {
			command, exists := commandsMap[cleanedWords[0]]
			if exists {
				var err error
				if (cleanedWords[0] == "explore" || cleanedWords[0] == "catch" || cleanedWords[0] == "inspect"|| cleanedWords[0] == "pokedex") && len(cleanedWords) > 1 {
					err = command.callback(&configPTR, cachePtr, cleanedWords[1], pokedex)
				} else if cleanedWords[0] == "explore" || cleanedWords[0] == "catch" {
					fmt.Println("Error: missing pokemon or location argument.")
				} else {
					err = command.callback(&configPTR, cachePtr, "", pokedex)
				}

				if err != nil {
					fmt.Printf("Error occurred: %v\n", err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}

		fmt.Print("Pokedex > ")
	}

	// Detect and report any error that happened during input scanning.
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
