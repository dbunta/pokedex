package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	pokecache "github.com/dbunta/pokedex/internal"
)

var cache pokecache.Cache

func startRepl() {

	fmt.Print("Welcome to the Pokedex!\n")
	config := config{next: "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20", previous: ""}

	duration := 10 * time.Second
	cache = pokecache.NewCache(duration)

	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		command := cleanInput(scanner.Text())

		if len(command) == 0 {
			os.Exit(0)
		}

		if value, ok := commands[command[0]]; ok {
			var arg string
			if len(command) > 1 {
				arg = command[1]
			}
			if err := value.callback(&config, arg); err != nil {
				fmt.Print("\n\n----error----\n\n")
				fmt.Print(err)
			}
		} else {
			fmt.Print("Unknown command\n")
		}

	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, string) error
}

type config struct {
	next         string
	previous     string
	locationArea string
}

// type locationArea struct {
// 	name string
// 	url  string
// }

type locationAreaResult struct {
	Count    int
	Next     string
	Previous string
	Results  []struct {
		Name string
		Url  string
	}
}

type locationArea struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
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
			description: "Gets next 20 map areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "map back",
			description: "Gets previous 20 map areas",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Gets specific map area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to catch a specific pokemon",
			callback:    commandCatch,
		},
	}
}

func commandExit(config *config, area string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(config *config, area string) error {
	fmt.Print("Welcome to the Pokedex!\n")
	fmt.Print("Usage:\n\n")

	for _, val := range getCommands() {
		fmt.Printf("%v: %v\n", val.name, val.description)
	}

	return nil
}

func commandMap(config *config, area string) error {

	cached, isCached := cache.Get(config.next)
	var body []byte

	if isCached {
		//fmt.Printf("fetched %v from CACHE\n", config.next)
		body = cached
	} else {
		//fmt.Printf("fetched %v from WEB\n", config.next)
		res, err := http.Get(config.next)
		if err != nil {
			return fmt.Errorf("error getting from pokeapi. %w", err)
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error getting from pokeapi. %w", err)
		}
		// fmt.Printf("added %v to CACHE\n", config.next)
		cache.Add(config.next, body)
	}

	var areas locationAreaResult
	err := json.Unmarshal(body, &areas)
	if err != nil {
		return fmt.Errorf("error getting from pokeapi. %w", err)
	}

	config.previous = areas.Previous
	config.next = areas.Next

	// fmt.Printf("Next: %v\n", config.next)
	fmt.Printf("Previous: %v\n", config.previous)
	for _, area := range areas.Results {
		fmt.Printf("%v\n", area.Name)
	}

	return nil
}

func commandMapBack(config *config, area string) error {
	if len(config.previous) <= 0 {
		config.previous = config.next
	}

	cached, isCached := cache.Get(config.previous)
	var body []byte
	if isCached {
		// fmt.Printf("fetched %v from CACHE\n", config.previous)
		body = cached
	} else {
		// fmt.Printf("fetched %v from WEB\n", config.previous)
		res, err := http.Get(config.previous)
		if err != nil {
			return fmt.Errorf("error getting from pokeapi. %w", err)
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error getting from pokeapi. %w", err)
		}
		// fmt.Printf("added %v to CACHE\n", config.previous)
		cache.Add(config.previous, body)
	}

	var areas locationAreaResult
	err := json.Unmarshal(body, &areas)
	if err != nil {
		return fmt.Errorf("error getting from pokeapi. %w", err)
	}

	if len(areas.Previous) > 0 {
		config.previous = areas.Previous
	}
	config.next = areas.Next

	// fmt.Printf("Next: %v\n", config.next)
	// fmt.Printf("Previous: %v\n", config.previous)
	for _, area := range areas.Results {
		fmt.Printf("%v\n", area.Name)
	}

	return nil
}

func commandExplore(config *config, area string) error {
	fmt.Printf("Exploring %v...\n", area)
	url := "https://pokeapi.co/api/v2/location-area/" + area + "/"

	cached, isCached := cache.Get(url)
	var body []byte
	if isCached {
		//fmt.Printf("\nfetched %v from CACHE\n", config.previous)
		body = cached
	} else {
		//fmt.Printf("\nfetched %v from WEB\n", config.previous)
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("error getting from pokeapi. %w", err)
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error getting from pokeapi. %w", err)
		}

		cache.Add(url, body)
	}

	var locationArea locationArea
	err := json.Unmarshal(body, &locationArea)
	if err != nil {
		return fmt.Errorf("error getting from pokeapi. %w", err)
	}

	for _, pe := range locationArea.PokemonEncounters {
		fmt.Printf("- %v\n", pe.Pokemon.Name)
	}

	return nil
}

func commandCatch(config *config, pokemon string) error {
	test := rand.Intn(5)
	fmt.Printf("\n%v\n", test)
	if test > 2 {
		fmt.Print("\ncaught it\n")
		return nil
	}
	fmt.Print("\nmissed it\n")
	return nil
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
