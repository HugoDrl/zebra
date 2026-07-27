package parser

import (
	"os"
)

func readFile(filename string) ([]byte, error) {
	c, err := os.ReadFile(filename)
	if err != nil {
		return nil, &FileError{File: filename, Err: err}
	}
	return c, nil
}

func splitLine(line string) []string {
	splitBySpace := true
	words := make([]string, 0)
	currWord := ""
	for _, letter := range line {
		if letter == '"' {
			splitBySpace = !splitBySpace
		}
		if letter == ' ' && splitBySpace {
			words = append(words, currWord)
			currWord = ""
			continue
		}
		currWord += string(letter)
	}
	if currWord != "" {
		words = append(words, currWord)
	}
	return words
}
