package main

import (
	"encoding/json"
	"fmt" // Package for formatted I/O (input/output)
	"io"
	"log"
	"net/http"
	"os" // Package for operating system functionalities, like exiting the program¨
	"github.com/Marcus-Gustafsson/pokedexCLI/internal"
	"math/rand"
	"strings"
	"github.com/fatih/color"
)

// cliCommand represents a command available in the CLI interface.
type cliCommand struct {
	name        string              // The name of the command (e.g., "help")
	description string              // A short description of this command
	callback    func(*config, *internal.Cache, string, map[string]pokemonDetails) error // The function executed when this command is invoked
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

type locationAreaDetails struct {
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type pokemonDetails struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order        any `json:"order"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []struct {
		Abilities []struct {
			Ability  any  `json:"ability"`
			IsHidden bool `json:"is_hidden"`
			Slot     int  `json:"slot"`
		} `json:"abilities"`
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
	} `json:"past_abilities"`
	PastTypes []any `json:"past_types"`
	Species   struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       string `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  string `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      string `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale string `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default"`
				FrontFemale  any    `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      string `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string `json:"back_default"`
				BackFemale       string `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      string `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           string `json:"back_default"`
					BackShiny             string `json:"back_shiny"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					BackTransparent       string `json:"back_transparent"`
					FrontDefault          string `json:"front_default"`
					FrontShiny            string `json:"front_shiny"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
					FrontTransparent      string `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       string `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  string `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      string `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale string `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  string `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

const maxBaseExp = 300 // max base exp for calculating chance to capture pokemon! (mew.base_experience = 270 exp, it was used as a threshold for the max base exp)


// commandsMap maps command names to their cliCommand handler definitions.
// It is initialized in the init() function to resolve dependency cycles.
var commandsMap map[string]cliCommand

// init initializes the commandsMap with all supported commands.
// This resolves the initialization cycle between commandsMap and handler definitions.
func init() {
	commandsMap = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Show all available commands with a brief description of each.",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Quit the Pokedex application immediately.",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "List the next page of Pokémon location areas you can explore.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "List the previous page of Pokémon location areas.",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Show all Pokémon that can be encountered in a specified location area.",
			callback:    explore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a Pokémon by name and add it to your Pokedex if successful.",
			callback:    catch,
		},
		"inspect": {
			name:        "inspect",
			description: "View detailed stats and information about a Pokémon you have caught.",
			callback:    inspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Display a list of all Pokémon you have successfully caught.",
			callback:    pokedex,
		},
	}
}

// commandExit terminates the CLI Pokedex application immediately.
// It now prints the goodbye message in yellow for extra flair!
func commandExit(configPtr *config, cachePtr *internal.Cache, location string, pokedex map[string]pokemonDetails) error {
    // Bright yellow bold goodbye for a positive, friendly signoff
    color.New(color.FgHiYellow, color.Bold).Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil // Unreachable, but required
}

// commandHelp prints information about all available CLI commands.
// It lists each command with its name and description.
func commandHelp(configPtr *config, cachePtr *internal.Cache, location string, pokedex map[string]pokemonDetails) error {
    color.New(color.FgCyan, color.Bold).Println("Welcome to the Pokedex!")
    fmt.Println("Usage:")
	fmt.Println()
    for _, cliCommand := range commandsMap {
        // Command name in bold yellow, description in white
        color.New(color.FgHiYellow, color.Bold).Printf("%v: ", cliCommand.name)
        color.White("%v\n", cliCommand.description)
    }
    return nil
}

// commandMap fetches and displays a paginated list of location areas from the PokeAPI.
// It uses a cache to avoid unnecessary HTTP requests for previously seen pages.
func commandMap(configPtr *config, cachePtr *internal.Cache, location string, pokedex map[string]pokemonDetails) error {
    var url string

    // Determine which URL to fetch: next page or default start page
    if configPtr.Next == "" {
        url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
    } else {
        url = configPtr.Next
    }

    // Try to get the response data from the cache first.
    val, ok := cachePtr.Get(url)
    if !ok {
        // Not in cache! Make HTTP request to fetch data from the API
        fmt.Println("DBG: (from API)")
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
    // Highlight each location name in green!
    for _, result := range locations.Results {
        color.New(color.FgHiGreen, color.Bold).Printf("%v\n", result.Name)
    }

    return nil
}

// commandMap shows the next page (or start) of Pokémon locations using the PokeAPI.
// Results are cached for efficiency.
func commandMapB(configPtr *config, cachePtr *internal.Cache, location string, pokedex map[string]pokemonDetails) error {

	var url string

	if configPtr.Previous == nil {
		color.New(color.FgHiBlack).Println("You're on the first page...")
		return nil
	} else {
		url = *configPtr.Previous
	}

    val, ok := cachePtr.Get(url)
    if !ok {
        // Not in cache! Make HTTP request to fetch data from the API
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
		color.New(color.FgHiGreen, color.Bold).Printf("%v\n", result.Name)
	}

	return nil
}


// commandMapB shows the previous page of Pokémon locations (or warns if on the first page) using the PokeAPI.
func explore(configPtr *config, cachePtr *internal.Cache, areaName string, pokedex map[string]pokemonDetails) error {
    // Construct the API URL for the provided area name.
    url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v/", areaName)


    // Try to get the response data from the cache first.
    val, ok := cachePtr.Get(url)
    if !ok {
        // Not in cache! Make HTTP request to fetch data from the API
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
    
	// Unmarshal HTTP response to extract Pokémon encounters from the JSON
    var areaDetails locationAreaDetails
    err := json.Unmarshal(val, &areaDetails)
    if err != nil {
        return err
    }
	color.New(color.FgCyan, color.Bold).Printf(
    "You venture into %s...\nThese wild Pokémon can be found here:\n",
    areaName,
	)
    // Each wild Pokémon in magenta and bold
    for _, result := range areaDetails.PokemonEncounters {
        color.New(color.FgHiMagenta, color.Bold).Printf(" - %v\n", result.Pokemon.Name)
    }
	fmt.Println()

    return nil
}


// catch attempts to catch a Pokémon by name, using a probability based on base experience.
// If caught, adds the Pokémon to the user's Pokedex.
func catch(configPtr *config, cachePtr *internal.Cache, pokemonName string, pokedex map[string]pokemonDetails) error {

    // Message indicating which pokemon we are trying to catch
	pokemonName = strings.ToLower(pokemonName)
	color.New(color.Bold).Printf("Throwing a Pokéball at %v...\n", pokemonName)

	
	// Construct the API URL for the provided area name.
    url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%v/", pokemonName)

    // Try to get the response data from the cache first.
    val, ok := cachePtr.Get(url)

    if !ok {
        // Not in cache! Make HTTP request to fetch data from the API
        fmt.Println("DBG: (from API)")
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
    
	// Unmarshal HTTP response to extract Pokémon encounters from the JSON
    var pokemon pokemonDetails
    err := json.Unmarshal(val, &pokemon)
    if err != nil {
        return err
    }

	chance := 1 - (float64(pokemon.BaseExperience) / float64(maxBaseExp))
	if chance < 0 {
		chance = 0.01
	}
	if chance > 1 {
		chance = 0.99
	}
	randomFloat := rand.Float64()
	if randomFloat < chance {
		pokedex[pokemonName] = pokemon
		color.New(color.FgHiGreen, color.Bold).Printf("%v was caught!\n", pokemonName)
        color.New(color.FgCyan).Println("You may now inspect it with the inspect command.")
	} else {
		color.New(color.FgHiRed, color.Bold).Println("Missed catch!")
	}
	fmt.Println()

    return nil
}



// inspect displays detailed information about a caught Pokémon.
// If the user hasn't caught this Pokémon yet, prints a message.
func inspect(configPtr *config, cachePtr *internal.Cache, pokemonName string, pokedex map[string]pokemonDetails) error {
    foundPokemon, ok := pokedex[pokemonName]
    if ok {
        // Name header
        color.New(color.FgHiYellow, color.Bold).Printf("Name: %v\n", foundPokemon.Name)
        color.New(color.Bold).Printf("Height: %v\nWeight: %v\n", foundPokemon.Height, foundPokemon.Weight)

        color.New(color.FgCyan, color.Bold).Println("Stats:")
        for _, value := range foundPokemon.Stats {
            statColor := color.New(color.Bold)
            switch value.Stat.Name {
            case "hp":
                statColor.Add(color.FgHiGreen)
            case "attack":
                statColor.Add(color.FgHiRed)
            case "defense":
                statColor.Add(color.FgBlue)
            case "special-attack":
                statColor.Add(color.FgHiMagenta)
            case "special-defense":
                statColor.Add(color.FgHiCyan)
            case "speed":
                statColor.Add(color.FgHiWhite)
            default:
                statColor.Add(color.FgWhite)
            }
            statColor.Printf("  - %v: %v\n", value.Stat.Name, value.BaseStat)
        }

        color.New(color.FgCyan, color.Bold).Println("Types:")
        for _, value := range foundPokemon.Types {
            typeName := value.Type.Name
            typeColor := color.New(color.Bold)
            switch typeName {
            case "fire":
                typeColor.Add(color.FgHiRed)
            case "water":
                typeColor.Add(color.FgHiCyan)
            case "grass":
                typeColor.Add(color.FgHiGreen)
            case "electric":
                typeColor.Add(color.FgHiYellow)
            case "psychic":
                typeColor.Add(color.FgMagenta)
            case "bug":
                typeColor.Add(color.FgGreen)
            case "normal":
                typeColor.Add(color.FgWhite)
            case "fighting":
                typeColor.Add(color.FgRed)
            case "poison":
                typeColor.Add(color.FgMagenta)
            case "ground":
                typeColor.Add(color.FgYellow)
            case "flying":
                typeColor.Add(color.FgHiBlue)
            case "rock":
                typeColor.Add(color.FgHiWhite)
            case "ghost":
                typeColor.Add(color.FgHiMagenta)
            case "ice":
                typeColor.Add(color.FgCyan)
            case "dragon":
                typeColor.Add(color.FgBlue)
            case "dark":
                typeColor.Add(color.FgBlack)
            case "steel":
                typeColor.Add(color.FgHiWhite)
            case "fairy":
                typeColor.Add(color.FgHiMagenta)
            default:
                typeColor.Add(color.FgCyan)
            }
            typeColor.Printf("  - %v\n", typeName)
        }
    } else {
        color.New(color.FgHiRed, color.Bold).Printf("You have not yet caught %v\n", pokemonName)
    }
    return nil
}

// pokedex lists all caught Pokémon names in the user's personal Pokedex.
func pokedex(configPtr *config, cachePtr *internal.Cache, pokemonName string, pokedex map[string]pokemonDetails) error {
    if len(pokedex) > 0 {
        color.New(color.FgCyan, color.Bold).Println("Your Pokedex:")
        for key, details := range pokedex {
            typeColor := color.New(color.Bold)
            if len(details.Types) > 0 {
                typeName := details.Types[0].Type.Name
                switch typeName {
                case "fire":
                    typeColor.Add(color.FgHiRed)
                case "water":
                    typeColor.Add(color.FgHiCyan)
                case "grass":
                    typeColor.Add(color.FgHiGreen)
                case "electric":
                    typeColor.Add(color.FgHiYellow)
                case "psychic":
                    typeColor.Add(color.FgMagenta)
                case "bug":
                    typeColor.Add(color.FgGreen)
                case "normal":
                    typeColor.Add(color.FgWhite)
                case "fighting":
                    typeColor.Add(color.FgRed)
                case "poison":
                    typeColor.Add(color.FgMagenta)
                case "ground":
                    typeColor.Add(color.FgYellow)
                case "flying":
                    typeColor.Add(color.FgHiBlue)
                case "rock":
                    typeColor.Add(color.FgHiWhite)
                case "ghost":
                    typeColor.Add(color.FgHiMagenta)
                case "ice":
                    typeColor.Add(color.FgCyan)
                case "dragon":
                    typeColor.Add(color.FgBlue)
                case "dark":
                    typeColor.Add(color.FgHiBlack)
                case "steel":
                    typeColor.Add(color.FgHiWhite)
                case "fairy":
                    typeColor.Add(color.FgHiMagenta)
                default:
                    typeColor.Add(color.FgCyan)
                }
            } else {
                typeColor.Add(color.FgHiGreen) // fallback for unknown type
            }
			color.New(color.Bold).Print("- ")
            typeColor.Printf("%v\n", key)
        }
    } else {
        color.New(color.FgHiMagenta, color.Bold).Println("No Pokémon in the Pokedex yet... Gotta catch 'em all!!")
    }
    return nil
}

