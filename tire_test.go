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

type PutValue struct {
	key            string
	value          interface{}
	keyIndex       int
	expectedSize   int
	expectedNewKey bool
	errorExpected  bool
}

func TestPut(t *testing.T) {
	var tests = []struct {
		params []PutValue
	}{
		{[]PutValue{ // two keys
			PutValue{"www.test.com", 1, 0, 1, true, false},
			PutValue{"www.example.com", 1, 0, 2, true, false},
		}},
		{[]PutValue{ // single key
			PutValue{"www.example.com", 1, 0, 1, true, false},
		}},
		{[]PutValue{ //key index > len(key)
			PutValue{"www.example.com", 1, 100, 1, true, true},
		}},
		{[]PutValue{ //key index < 0
			PutValue{"www.example.com", 1, -1, 1, true, true},
		}},
		{[]PutValue{ //keyIndex == len(key)
			PutValue{"www.example.com", 1, len("www.example.com"), 0, true, false},
		}},
		{[]PutValue{ //keyIndex == len(key), update value
			PutValue{"www.example.com", 1, len("www.example.com"), 0, true, false},
			PutValue{"www.example.com", 2, len("www.example.com"), 0, false, false},
		}},
		{[]PutValue{ //multiple values
			PutValue{"www.test.com", 1, 0, 1, true, false},
			PutValue{"www.example.com", 2, 0, 2, true, false},
		}},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.params)
		t.Run(testname, func(t *testing.T) {
			node, _ := newTrieNode("", nil, 0)
			for _, param := range tt.params {
				newKey, err := node.put(param.key, param.value, param.keyIndex)
				if !param.errorExpected && err != nil {
					t.Errorf("Error was not expected")
				}
				if param.errorExpected && err == nil {
					t.Errorf("Error was expected")
				}
				if !param.errorExpected && node.size != param.expectedSize {
					t.Errorf("Expected size %d go %d", param.expectedSize, node.size)
				}
				if !param.errorExpected && newKey != param.expectedNewKey {
					t.Errorf("Expected newKey %t go %t", param.expectedNewKey, newKey)
				}
			}
		})
	}
}
