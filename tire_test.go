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
				t.Errorf("Error was expected")
			}
			if !tt.errorExpected && node.size != tt.expectedSize {
				t.Errorf("Expected size %d go %d", tt.expectedSize, node.size)
			}
		})
	}
}

func TestPut(t *testing.T) {
	var tests = []struct {
		keyValues     map[string]interface{}
		expectedSize  int
		errorExpected bool
	}{
		{map[string]interface{}{
			"www.test.com": 1,
		},
			1, false},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.keyValues)
		t.Run(testname, func(t *testing.T) {
			node, err := newTrieNode("", nil, 0)
			for key, value := range tt.keyValues {
				node.put(key, value, 0)
			}
			if !tt.errorExpected && err != nil {
				t.Errorf("Error was not expected")
			}
			if tt.errorExpected && err == nil {
				t.Errorf("Error was expected")
			}
			if !tt.errorExpected && node.size != tt.expectedSize {
				t.Errorf("Expected size %d go %d", tt.expectedSize, node.size)
			}
		})
	}
}
