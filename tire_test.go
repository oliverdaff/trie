package trie

import (
	"fmt"
	"testing"
)

func TestNewTrieNode(t *testing.T) {
	var tests = []struct {
		key           string
		value         interface{}
		keyIndex      int
		errorExpected bool
	}{
		{"www.testing.com", nil, 0, false},
		{"www.testing.com", nil, 100, true},
		{"www.testing.com", nil, -1, true},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.key)
		t.Run(testname, func(t *testing.T) {
			_, err := newTrieNode(tt.key, tt.value, tt.keyIndex)
			if !tt.errorExpected && err != nil {
				t.Errorf("Error was not expected")
			}
			if tt.errorExpected && err == nil {
				t.Errorf("Error was not expected")
			}
		})
	}
}
