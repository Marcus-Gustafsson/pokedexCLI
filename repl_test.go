package main

import (
	"testing"
)

// TestCleanInput tests the cleanInput function with various input scenarios
// It uses table-driven testing to check multiple test cases efficiently
func TestCleanInput(t *testing.T) {
	// Define test cases using an anonymous struct slice
	// Each test case has an input string and expected output slice
	cases := []struct {
		input    string   // The raw input string to test
		expected []string // The expected cleaned output
	}{
		{
			// Test case 1: Input with extra whitespace around words
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			// Test case 2: Mixed case Pokemon names with extra whitespace
			input:    "  Charmander Bulbasaur PIKACHU  ",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
	}

	// Iterate through each test case
	for _, c := range cases {
		// Call the function we're testing with the current test input
		actual := cleanInput(c.input)
		
		// First check: Verify the length of the result matches expected length
		if len(actual) != len(c.expected) {
			t.Errorf("Expected length: %v, Got length: %v", len(actual), len(c.expected))
		}
		
		// Second check: Compare each word in the result with expected words
		for i := range actual {
			word := actual[i]           // The actual word at position i
			expectedWord := c.expected[i] // The expected word at position i

			// If words don't match, report the error
			if word != expectedWord {
				t.Errorf("Expected word: %s, Got word: %s", expectedWord, word)
			}
		}
	}
}


