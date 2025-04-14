package main

import (
  "fmt"
  "strings"
  "bufio"
  "os"
)

func main() {
  for {
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Print("Pokedex > ")
    _ = scanner.Scan()
    text := scanner.Text()
    words := strings.Split(strings.ToLower(text), " ")
    fmt.Printf("Your command was: %v\n", words[0])
  }
}

func cleanInput(text string) []string {
  words := strings.Split(strings.Trim(strings.ToLower(text), " "), " ")
  var finalWords []string
  for _, word := range words {
    word = strings.Replace(word, " ", "", -1)
    fmt.Print(word)
    if len(word) > 0 {
      finalWords = append(finalWords, word)
    }
  }
  return finalWords
}
