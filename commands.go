package main

import (
	"fmt" // Package for formatted I/O (input/output)
	"os"  // Package for operating system functionalities, like exiting the program
)

// cliCommand represents a single command that can be run in our command-line interface (CLI).
type cliCommand struct {
	name        string
	description string
	callback    func() error // A function that takes no arguments and returns an error
}

// commandsMap is a global registry of all supported CLI commands.
// It maps the command name (a string) to its corresponding cliCommand struct.
// Its initialization (populating it with actual commands) happens in the init() function.
var commandsMap map[string]cliCommand

// init() is a special Go function that runs automatically before the main() function.
// We use it here to initialize our 'commandsMap'. This is crucial for solving
// the "initialization cycle" problem. If we tried to initialize commandsMap directly
// with 'commandHelp' (which itself refers to commandsMap to print all commands),
// Go would detect a circular dependency during compilation.
// By populating the map in init(), both commandExit and commandHelp functions
// are fully defined and ready to be referenced by the map, resolving the cycle.
func init() {
	commandsMap = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp, // Assigning the commandHelp function as the callback
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit, // Assigning the commandExit function as the callback
		},
	}
}

// commandExit handles the "exit" command.
func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil // This line is technically unreachable due to os.Exit(0), but required by the function signature
}

// commandHelp handles the "help" command.
// It prints a welcome message, then iterates through the 'commandsMap'
// to dynamically display the name and description of each registered command.
func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println() // Prints an empty line for better formatting
	for _, cliCommand := range commandsMap {
		// Prints each command's name and description in the format "name: description"
		fmt.Printf("%v: %v\n", cliCommand.name, cliCommand.description)
	}
	return nil // needed for the function signature
}