package parser

import (
	"slices"
	"testing"
)


func TestSplitLine(t *testing.T) {
	tests := map[string]struct{
		input string
		expected []string
	}{
		"one word should return the word": {
			input: "hello",
			expected: []string{"hello"},
		},
		"one sentence should split each words": {
			input: "this is a sentence",
			expected: []string{"this", "is", "a", "sentence"},
		},
		"sentence with quotes should keep words in quotes": {
			input: "this is a sentence \"with several words in quotes\"",
			expected: []string{"this", "is", "a", "sentence", "with several words in quotes"},
		},
		"words with quotes in it should not do anything": {
			input: "this is a sen\"tence",
			expected: []string{"this", "is", "a", "sen\"tence"},
		},
	}

	for name, test := range tests {
		output := splitLine(test.input)
		if !slices.Equal(output, test.expected) {
			t.Fatalf("%s: expected %v, got %v", name, test.expected, output)
		}
	}
}
