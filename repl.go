package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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
			if err := value.callback(&config); err != nil {
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
	callback    func(*config) error
}

type config struct {
	next     string
	previous string
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
	}
}

func commandExit(config *config) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(config *config) error {
	fmt.Print("Welcome to the Pokedex!\n")
	fmt.Print("Usage:\n\n")

	for _, val := range getCommands() {
		fmt.Printf("%v: %v\n", val.name, val.description)
	}

	return nil
}

func commandMap(config *config) error {

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

func commandMapBack(config *config) error {
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

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
