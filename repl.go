package main

import (
	"strings"
)

// cleanInput takes a raw input string and returns a slice of cleaned words
// It converts to lowercase and splits on whitespace, removing empty strings
func cleanInput(text string) []string {
	// Convert to lowercase for case-insensitive commands
	lower := strings.ToLower(text)
	// Split on whitespace and filter out empty strings
	words := strings.Fields(lower)
	return words
}
