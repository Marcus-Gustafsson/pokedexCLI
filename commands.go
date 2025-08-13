package main

import (
	"encoding/json"
	"fmt" // Package for formatted I/O (input/output)
	"io"
	"log"
	"net/http"
	"os" // Package for operating system functionalities, like exiting the program¨
	"github.com/Marcus-Gustafsson/pokedexCLI/internal"
)

// cliCommand represents a command available in the CLI interface.
type cliCommand struct {
	name        string              // The name of the command (e.g., "help")
	description string              // A short description of this command
	callback    func(*config, *internal.Cache) error // The function executed when this command is invoked
}

// locations holds the results from the PokeAPI location-area endpoint.
type locations struct {
	Results []struct {
		Name string `json:"name"` // Name of the location area
		URL  string `json:"url"`  // API URL for details about this location area
	} `json:"results"`
}

// config stores pagination URLs and result count for navigating paginated PokeAPI responses.
type config struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`     // URL for the next set of results, or "" if none
	Previous *string `json:"previous"` // URL for the previous set, or nil if on the first page
}

// commandsMap maps command names to their cliCommand handler definitions.
// It is initialized in the init() function to resolve dependency cycles.
var commandsMap map[string]cliCommand

// init initializes the commandsMap with all supported commands.
// This resolves the initialization cycle between commandsMap and handler definitions.
func init() {
	commandsMap = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
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

// commandExit exits the program immediately.
// Required by the CLI to terminate gracefully.
func commandExit(configPtr *config, cachePtr *internal.Cache) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil // This line is technically unreachable due to os.Exit(0), but required by the function signature
}

// commandHelp prints information about all available CLI commands.
// It lists each command with its name and description.
func commandHelp(configPtr *config, cachePtr *internal.Cache) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cliCommand := range commandsMap {
		fmt.Printf("%v: %v\n", cliCommand.name, cliCommand.description)
	}
	return nil // needed for the function signature
}

// commandMap fetches and displays a paginated list of location areas from the PokeAPI.
// It uses a cache to avoid unnecessary HTTP requests for previously seen pages.
func commandMap(configPtr *config, cachePtr *internal.Cache) error {
    var url string

    // Determine which URL to fetch: next page or default start page
    if configPtr.Next == "" {
        url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
    } else {
        url = configPtr.Next
    }

    // Try to get the response data from the cache first.
	fmt.Println("Looking up URL:", url)
    val, ok := cachePtr.Get(url)
    if ok {
        // Cached response found; skip the network call, saves time and bandwidth
        fmt.Println("(from cache)")
    } else {
        // Not in cache! Make HTTP request to fetch data from the API
        fmt.Println("(from API)")
        res, err := http.Get(url)
        if err != nil {
            return err
        }

        defer res.Body.Close() // Always close response body when done

        val, err = io.ReadAll(res.Body)
        if err != nil {
            return err
        }

        if res.StatusCode > 299 {
            // Log if response is an error
            log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, val)
        }

        // Store the raw byte response in the cache for next time
        cachePtr.Add(url, val)
    }
    
    var locations locations
    err := json.Unmarshal(val, &locations)
    if err != nil {
        return err
    }

    err = json.Unmarshal(val, configPtr)
    if err != nil {
        return err
    }

    // Print the names of all locations in the page
    for _, result := range locations.Results {
        fmt.Printf("%v\n", result.Name)
    }

    return nil
}

// commandMapB (map back) fetches and displays the previous 20 location areas from the PokeAPI.
// If already at the first page, it informs the user.
func commandMapB(configPtr *config, cachePtr *internal.Cache) error {

	var url string

	if configPtr.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	} else {
		url = *configPtr.Previous
	}

	// Try to get the response data from the cache first.
	fmt.Println("Looking up URL:", url)
    val, ok := cachePtr.Get(url)
    if ok {
        // Cached response found; skip the network call, saves time and bandwidth
        fmt.Println("(from cache)")
    } else {
        // Not in cache! Make HTTP request to fetch data from the API
        fmt.Println("(from API)")
        res, err := http.Get(url)
        if err != nil {
            return err
        }

        defer res.Body.Close() // Always close response body when done

        val, err = io.ReadAll(res.Body)
        if err != nil {
            return err
        }

        if res.StatusCode > 299 {
            // Log if response is an error
            log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, val)
        }

        // Store the raw byte response in the cache for next time
        cachePtr.Add(url, val)
    }

	locations := locations{}

	err := json.Unmarshal(val, &locations)
	if err != nil {
		return err
	}

	err = json.Unmarshal(val, configPtr)
	if err != nil {
		return err
	}

	for _, result := range locations.Results {
		fmt.Printf("%v\n", result.Name)
	}

	return nil
}
