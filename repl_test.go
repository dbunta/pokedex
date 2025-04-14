package main

import "testing"

func TestCleanInput(t *testing.T) {
  cases := []struct {
    input string
    expected []string
  }{
    {
      input: "   hello world   ",
      expected: []string{"hello", "world"},
    },
    {
      input: "APPLE    jAcKs",
      expected: []string{"apple", "jacks"},
    },
    {
      input: "   APPLE jAcKs   ",
      expected: []string{"apple", "jacks"},
    },
  }

  for _, c := range cases {
    actual := cleanInput(c.input)
    // Check the length of the actual slice against the expected slice
    // if they don't match, use t.Errorf to print an error message
    // and fail the test

    if len(actual) != len(c.expected) {
      t.Errorf("Actual words count and expected words count do not match. Expected: %v, Actual: %v", len(c.expected), len(actual))
        return
    }

    for i := range actual {
      word := actual[i]
      expectedWord := c.expected[i]
      // Check each word in the slice
      // if they don't match, use t.Errorf to print an error message
      // and fail the test
      if word != expectedWord {
        t.Errorf("Expected word does not match actual. Expected: %v, Actual: %v", expectedWord, word)
      }
    }
  }
}
