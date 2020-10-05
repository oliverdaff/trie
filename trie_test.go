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
		{"www.testing.com", 1, 1},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.key)
		t.Run(testname, func(t *testing.T) {
			node, _ := newTrieNode(tt.key, tt.value)
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

type DeleteValue struct {
	key            string
	deleted, empty bool
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
		{[]PutValue{ //update value
			PutValue{"www.test.com", 1, 1, true},
			PutValue{"www.test.com", 2, 1, false},
		}},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.params)
		t.Run(testname, func(t *testing.T) {
			node, _ := newTrieNode("", nil)
			for _, param := range tt.params {
				newKey, _ := node.put(param.key, param.value)
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

func TestGetNode(t *testing.T) {
	var tests = []struct {
		params    []PutValue
		query     string
		nodeFound bool
	}{
		{[]PutValue{ // node found
			PutValue{"www.test.com", 1, 1, true},
		}, "www.test.com", true},
		{[]PutValue{ // node not found
			PutValue{"www.example.com", 1, 1, true},
		}, "www.test.com", false},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.params)
		t.Run(testname, func(t *testing.T) {
			node, _ := newTrieNode("", 1)
			for _, param := range tt.params {
				newKey, _ := node.put(param.key, param.value)
				if node.size != param.expectedSize {
					t.Errorf("Expected size %d go %d", param.expectedSize, node.size)
				}
				if newKey != param.expectedNewKey {
					t.Errorf("Expected newKey %t go %t", param.expectedNewKey, newKey)
				}
			}
			queryResult := node.getNode(tt.query)
			if tt.nodeFound && queryResult == nil {
				t.Errorf("Expected node for query %s got nil", tt.query)
			}
			if !tt.nodeFound && queryResult != nil {
				t.Errorf("Expected null node for query %s got %v", tt.query, queryResult)
			}
		})
	}
}

func TestContains(t *testing.T) {
	var tests = []struct {
		params   []PutValue
		query    string
		keyFound bool
	}{
		{[]PutValue{ // node found
			PutValue{"www.test.com", 1, 1, true},
		}, "www.test.com", true},
		{[]PutValue{ // node not found
			PutValue{"www.example.com", 1, 1, true},
		}, "www.test.com", false},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.params)
		t.Run(testname, func(t *testing.T) {
			node, _ := newTrieNode("", 1)
			for _, param := range tt.params {
				node.put(param.key, param.value)
			}
			keyFound := node.contains(tt.query)
			if tt.keyFound != keyFound {
				t.Errorf("Expected key found %t got %t", tt.keyFound, keyFound)
			}
		})
	}
}

func TestGet(t *testing.T) {
	var tests = []struct {
		params []PutValue
		query  string
		value  interface{}
	}{
		{[]PutValue{ // Get value
			PutValue{"www.test.com", 1, 1, true},
		}, "www.test.com", 1},
		{[]PutValue{ // Get Value
			PutValue{"www.example.com", 2, 1, true},
		}, "www.example.com", 2},
		{[]PutValue{ // Value not found
			PutValue{"www.example.com", 2, 1, true},
		}, "www.test.com", nil},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.params)
		t.Run(testname, func(t *testing.T) {
			node, _ := newTrieNode("", 1)
			for _, param := range tt.params {
				node.put(param.key, param.value)
			}
			value := node.get(tt.query)
			if tt.value != value {
				t.Errorf("Expected value found %d got %d", tt.value, value)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	var tests = []struct {
		params    []PutValue
		deletions []DeleteValue
	}{
		{[]PutValue{ // Delete value
			PutValue{"www.test.com", 1, 1, true},
		},
			[]DeleteValue{
				DeleteValue{"www.test.com", true, true},
				DeleteValue{"www.test.com", false, false},
			},
		},
		{[]PutValue{ // Delete partial key
			PutValue{"www.test.com", 1, 1, true},
		},
			[]DeleteValue{
				DeleteValue{"www.test", false, false},
			},
		},
		{[]PutValue{ // Delete shared key
			PutValue{"www.test.com", 1, 1, true},
			PutValue{"www.example.com", 1, 1, true},
		},
			[]DeleteValue{
				DeleteValue{"www.test.com", true, false},
				DeleteValue{"www.example.com", true, true},
			},
		},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.params)
		t.Run(testname, func(t *testing.T) {
			node, err := newTrieNode("", nil)
			if err != nil {
				t.Errorf("Error not expected")
			}
			for _, param := range tt.params {
				node.put(param.key, param.value)
			}
			for _, deletion := range tt.deletions {
				deleted, empty := node.delete(deletion.key)
				if deletion.deleted != deleted {
					t.Errorf("Expected deleted %t got %t", deletion.deleted, deleted)
				}
				if deletion.empty != empty {
					t.Errorf("Expected empty %t got %t", deletion.empty, empty)
				}

			}
		})
	}
}

func TestLongestPrefix(t *testing.T) {
	var tests = []struct {
		params         []PutValue
		prefix         string
		expectedPrefix string
	}{
		{[]PutValue{
			PutValue{"www.test.com", 1, 1, true},
		},
			"www.test.com", "www.test.com",
		},
		{[]PutValue{
			PutValue{"www.test.com", 1, 1, true},
			PutValue{"www.test", 1, 1, true},
		},
			"www.test.co.uk", "www.test",
		},
		{[]PutValue{
			PutValue{"www.test.com", 1, 1, true},
			PutValue{"www.test", 1, 1, true},
			PutValue{"www", 1, 1, true},
		},
			"www.example.com", "www",
		},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.params)
		t.Run(testname, func(t *testing.T) {
			node, _ := newTrieNode("", 1)
			for _, param := range tt.params {
				node.put(param.key, param.value)
			}
			prefix := node.longestPrefixOf(tt.prefix, 0)
			if tt.expectedPrefix != prefix {
				t.Errorf("Expected prefix %s got %s", tt.expectedPrefix, prefix)
			}

		})
	}
}
