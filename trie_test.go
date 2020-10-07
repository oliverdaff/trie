package trie

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewTrie(t *testing.T) {
	trie := NewTrie()
	if trie == nil {
		t.Error("Trie must not be null")
	}
}

func TestTriePut(t *testing.T) {
	var empty int
	var tests = []struct {
		key           string
		value         interface{}
		errorExpected bool
	}{
		{"www.testing.com", 1, false},
		{"www.testing.com", nil, true},
		{"www.testing.com", empty, false},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.key)
		t.Run(testname, func(t *testing.T) {
			trie := NewTrie()
			err := trie.Put(tt.key, tt.value)
			if err == nil && tt.errorExpected {
				t.Errorf("Expected a error")
			}
			if err != nil && !tt.errorExpected {
				t.Errorf("Error not expected %s", err)
			}
		})
	}
}

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

func TestTrieKeysWithPrefix(t *testing.T) {
	var tests = []struct {
		keys            []string
		prefix          string
		lenReturnedKeys int
	}{
		{[]string{"www.test.com"}, "www", 1},
		{[]string{
			"www.test.com",
			"www.example.com",
		}, "www", 2},
		{[]string{
			"www.test.com",
			"www.example.com",
			"example.com",
		}, "www", 2},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.keys)
		t.Run(testname, func(t *testing.T) {
			trie := NewTrie()
			var empty int
			for _, key := range tt.keys {
				trie.Put(key, empty)
			}
			returnedKeys := make([]string, 0)
			for returned := range trie.KeysWithPrefix(tt.prefix) {
				returnedKeys = append(returnedKeys, returned)
			}
			if tt.lenReturnedKeys != len(returnedKeys) {
				t.Errorf("Expected %d keys got %d",
					tt.lenReturnedKeys, len(returnedKeys))
			}

		})
	}
}

func TestTrieKeys(t *testing.T) {
	var tests = []struct {
		keys []string
	}{
		{[]string{"www.test.com"}},
		{[]string{
			"www.test.com",
			"www.example.com",
		}},
		{[]string{
			"www.test.com",
			"www.example.com",
			"example.com",
		}},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.keys)
		t.Run(testname, func(t *testing.T) {
			trie := NewTrie()
			var empty int
			for i, key := range tt.keys {
				trie.Put(key, empty)
				if i+1 != trie.Size() {
					t.Errorf("Expected size %d got %d", i, trie.Size())
				}
			}
			returnedKeys := make([]string, 0)
			for returned := range trie.Keys() {
				returnedKeys = append(returnedKeys, returned)
			}
			if len(tt.keys) != len(returnedKeys) {
				t.Errorf("Expected %d key got %d : %s",
					len(tt.keys), len(returnedKeys), strings.Join(returnedKeys, ","))
			}

		})
	}
}

func TestTrieItems(t *testing.T) {
	var tests = []struct {
		keys []string
	}{
		{[]string{"www.test.com"}},
		{[]string{
			"www.test.com",
			"www.example.com",
		}},
		{[]string{
			"www.test.com",
			"www.example.com",
			"example.com",
		}},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.keys)
		t.Run(testname, func(t *testing.T) {
			trie := NewTrie()
			var empty int
			for i, key := range tt.keys {
				trie.Put(key, empty)
				if i+1 != trie.Size() {
					t.Errorf("Expected size %d got %d", i, trie.Size())
				}
			}
			returnedKeys := make([]NodeKeyValue, 0)
			for returned := range trie.Items() {
				returnedKeys = append(returnedKeys, returned)
			}
			if len(tt.keys) != len(returnedKeys) {
				t.Errorf("Expected %d key got %d ",
					len(tt.keys), len(returnedKeys))
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

func TestItems(t *testing.T) {
	var tests = []struct {
		params        []PutValue
		nodekeyValues []NodeKeyValue
	}{
		{[]PutValue{
			PutValue{"a", 1, 1, true},
		},
			[]NodeKeyValue{
				NodeKeyValue{"a", 1},
			},
		},
		{[]PutValue{
			PutValue{"b", 1, 1, true},
			PutValue{"a", 1, 1, true},
		},
			[]NodeKeyValue{
				NodeKeyValue{"a", 1},
				NodeKeyValue{"b", 1},
			},
		},
		{[]PutValue{
			PutValue{"ab", 1, 1, true},
			PutValue{"aa", 1, 1, true},
			PutValue{"a", 1, 1, true},
		},
			[]NodeKeyValue{
				NodeKeyValue{"a", 1},
				NodeKeyValue{"aa", 1},
				NodeKeyValue{"ab", 1},
			},
		},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.params)
		t.Run(testname, func(t *testing.T) {
			node, _ := newTrieNode("", nil)
			for _, param := range tt.params {
				node.put(param.key, param.value)
			}
			itemsChan := node.items(make([]byte, 0))
			items := make([]NodeKeyValue, 0)
			for item := range itemsChan {
				items = append(items, item)
			}

			for i := 0; i < len(tt.nodekeyValues); i++ {
				if items[i].key != tt.nodekeyValues[i].key {
					t.Errorf("Expected key %s got %s", tt.nodekeyValues[i].key, items[i].key)
				}
				if items[i].value != tt.nodekeyValues[i].value {
					t.Errorf("Expected value %s got %s", tt.nodekeyValues[i].value, items[i].value)
				}
			}

		})
	}
}
