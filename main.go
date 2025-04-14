package main

import (
  "fmt"
  "strings"
)

func main() {
  fmt.Print("Hello, World!")
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
