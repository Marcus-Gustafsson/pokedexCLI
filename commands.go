package main

import (
	"encoding/json"
	"fmt" // Package for formatted I/O (input/output)
	"io"
	"log"
	"net/http"
	"os" // Package for operating system functionalities, like exiting the program
)

// cliCommand represents a single command that can be run in our command-line interface (CLI).
type cliCommand struct {
	name        string
	description string
	callback    func(*config) error // A function that takes no arguments and returns an error
}

type locations struct {
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type config struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
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
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Lists the locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Lists the previous locations",
			callback:    commandMapB,
		},
	}
}

// commandExit handles the "exit" command.
func commandExit(configPTR *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil // This line is technically unreachable due to os.Exit(0), but required by the function signature
}

// commandHelp handles the "help" command.
// It prints a welcome message, then iterates through the 'commandsMap'
// to dynamically display the name and description of each registered command.
func commandHelp(configPTR *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println() // Prints an empty line for better formatting
	for _, cliCommand := range commandsMap {
		// Prints each command's name and description in the format "name: description"
		fmt.Printf("%v: %v\n", cliCommand.name, cliCommand.description)
	}
	return nil // needed for the function signature
}

func commandMap(configPTR *config) error {

	var url string

	if configPTR.Next == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		url = configPTR.Next
	}
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)

	defer res.Body.Close()

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

	locations := locations{}

	err = json.Unmarshal(body, &locations)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(body, configPTR)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\nDBG: configPTR.NEXT = %v, configPTR.PREVIOUS = %v\n", configPTR.Next, configPTR.Previous)

	for _, result := range locations.Results {
		fmt.Printf("%v\n", result.Name)
	}

	return nil
}

func commandMapB(configPTR *config) error {

	var url string

	if configPTR.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	} else {
		url = *configPTR.Previous
	}
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)

	defer res.Body.Close()

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

	locations := locations{}

	err = json.Unmarshal(body, &locations)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(body, configPTR)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\nDBG: configPTR.NEXT = %v, configPTR.PREVIOUS = %v\n", configPTR.Next, configPTR.Previous)

	for _, result := range locations.Results {
		fmt.Printf("%v\n", result.Name)
	}

	return nil
}
