package main

import (
	"bufio"
	"fmt"
	"os"
)

// main function starts the Pokedex REPL (Read-Eval-Print Loop)
func main() {

	// Create a new scanner to read from standard input (keyboard)
	scanner := bufio.NewScanner(os.Stdin)
	
	// Print the initial prompt
	fmt.Print("Pokedex > ")
	
	// Loop continues as long as there's input to scan (false when input ends (Ctrl+C or Ctrl+D))
	for scanner.Scan() {
		// Get the user's input text
		userInput := scanner.Text()
		
		// Clean and parse the input into words
		cleanedWords := cleanInput(userInput)
		
		// Check if we have at least one word (command)
		if len(cleanedWords) > 0 {
			// Print the first word as the command
			fmt.Printf("DBG: Your command was: %v \n", cleanedWords[0])
			command, exists := commandsMap[cleanedWords[0]]
			if exists{
				fmt.Printf("DBG: command: **%v** was found, running function...\n", cleanedWords[0])
				err := command.callback()
				if err != nil{
					fmt.Printf("Error occured: %v \n", err)
				}
			}else{
				fmt.Println("Unknown command")
			}
			
		}

		// Print the prompt again for the next command
		fmt.Print("Pokedex > ")
	}
	
	// Handle any scanning errors
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

}
