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
		expectedSize  int
		errorExpected bool
	}{
		{"www.testing.com", nil, 0, 1, false},
		{"www.testing.com", nil, len("www.testing.com"), 0, false},
		{"www.testing.com", nil, 100, 0, true},
		{"www.testing.com", nil, -1, 0, true},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.key)
		t.Run(testname, func(t *testing.T) {
			node, err := newTrieNode(tt.key, tt.value, tt.keyIndex)
			if !tt.errorExpected && err != nil {
				t.Errorf("Error was not expected")
			}
			if tt.errorExpected && err == nil {
				t.Errorf("Error was not expected")
			}
			if !tt.errorExpected && node.size != tt.expectedSize {
				t.Errorf("Expected size %d go %d", tt.expectedSize, node.size)
			}
		})
	}
}
