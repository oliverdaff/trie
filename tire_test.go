package trie

import (
	"fmt"
	"testing"
)

func TestNewTrieNode(t *testing.T) {
	var tests = []struct {
		key          string
		value        interface{}
		expectedSize int
	}{
		{"www.testing.com", nil, 1},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.key)
		t.Run(testname, func(t *testing.T) {
			node := newTrieNode(tt.key, tt.value)
			if node.size != tt.expectedSize {
				t.Errorf("Expected size %d go %d", tt.expectedSize, node.size)
			}
		})
	}
}

type PutValue struct {
	key            string
	value          interface{}
	expectedSize   int
	expectedNewKey bool
}

func TestPut(t *testing.T) {
	var tests = []struct {
		params []PutValue
	}{
		{[]PutValue{ // two keys
			PutValue{"www.test.com", 1, 1, true},
			PutValue{"www.example.com", 1, 2, true},
		}},
		{[]PutValue{ // single key
			PutValue{"www.example.com", 1, 1, true},
		}},
		{[]PutValue{ //multiple values
			PutValue{"www.test.com", 1, 1, true},
			PutValue{"www.example.com", 2, 2, true},
		}},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.params)
		t.Run(testname, func(t *testing.T) {
			node := newTrieNode("", nil)
			for _, param := range tt.params {
				newKey := node.put(param.key, param.value)
				if node.size != param.expectedSize {
					t.Errorf("Expected size %d go %d", param.expectedSize, node.size)
				}
				if newKey != param.expectedNewKey {
					t.Errorf("Expected newKey %t go %t", param.expectedNewKey, newKey)
				}
			}
		})
	}
}
