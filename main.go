package main

import (
  "fmt"
  "strings"
  "bufio"
  "os"
)


func main() {

  fmt.Print("Welcome to the Pokedex!\n")

  scanner := bufio.NewScanner(os.Stdin)
  commands := getCommands()
  for {
    scanner.Scan()
    command := cleanInput(scanner.Text())

    fmt.Print("Pokedex > ")

    if len(command) == 0 {
      os.Exit(0)
    }

    if value, ok := commands[command[0]]; ok {
      value.callback()
    } else {
      fmt.Print("Unknown command\n")
    }

  }
}

type cliCommand struct {
  name string
  description string
  callback func() error
}

func getCommands() map[string] cliCommand {
  return map[string] cliCommand {
    "exit": {
      name: "exit",
      description: "Exit the Pokedex",
      callback: commandExit,
    },
    "help": {
      name: "help",
      description: "Displays a help message",
      callback: commandHelp,
    },
  }
}


func commandExit() error {
  fmt.Printf("Closing the Pokedex... Goodbye!\n")
  os.Exit(0)
  return fmt.Errorf("exit command")
}

func commandHelp() error {
  fmt.Print("Welcome to the Pokedex!\n")
  fmt.Print("Usage:\n\n")

  for _, val := range getCommands() {
    fmt.Printf("%v: %v\n", val.name, val.description)
  }

  return fmt.Errorf("help command")
}

func cleanInput(text string) []string {
  words := strings.Split(strings.Trim(strings.ToLower(text), " "), "\n")
  var finalWords []string
  for _, word := range words {
    word = strings.Replace(word, " ", "", -1)
    if len(word) > 0 {
      finalWords = append(finalWords, word)
    }
  }
  return finalWords
}




