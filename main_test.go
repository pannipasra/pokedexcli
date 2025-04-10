package main

import (
	"testing"
)

func TestCleanInput(t *testing.T){
	// ?
	cases := []struct {
		input		string
		expected	[]string
	}{
		{
			input:	"  hello world   ",
			expected: []string{"hello", "world"},
		},
		// add more test cases here..
	}


	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		// Check the length of the actual slice against the expected slice
		if len(actual) != len(c.expected) {
			t.Errorf("Length mismatch for input '%s': got %d words, expected %d words", 
				c.input, len(actual), len(c.expected))
			continue
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("Word mismatch for input '%s': got '%s' at position %d, expected '%s'", c.input, word, i, expectedWord)
			}
		}
	}
}
