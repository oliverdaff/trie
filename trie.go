package trie

import (
	"errors"
	"sort"
)

// trieNode is a internal representation of a trie.
// Each node is root of its sub-trie. trieNode allows searching and adding new key-value pairs.
type trieNode struct {
	links map[byte]*trieNode
	size  int
	value interface{}
}

// newTrieNode creates a new trie
//
// - key: the possibly empty key to store in the trie
// - value: the value to be associated with the key
func newTrieNode(key string, value interface{}) (*trieNode, error) {
	if value == nil && len(key) != 0 {
		return nil, errors.New("nil value passed to newTrieNode with non null key")
	}
	links := make(map[byte]*trieNode)
	if len(key) == 0 {
		return &trieNode{
			size:  0,
			value: value,
			links: links,
		}, nil
	}
	node, err := newTrieNode(key[1:], value)
	if err != nil {
		return nil, err
	}
	links[key[0]] = node

	return &trieNode{
		size:  1,
		links: links,
	}, nil

}

// put stores a key value pair in the trie.
// Returns true unless the key was already in the trie and got updated.
func (ts *trieNode) put(key string, value interface{}) (bool, error) {
	if len(key) == 0 {
		isNewKey := ts.value == nil
		ts.value = value
		return isNewKey, nil
	}
	next := key[0]
	if nextNode, ok := ts.links[next]; ok {
		isNewKey, err := nextNode.put(key[1:], value)
		if err != nil {
			return false, err
		}
		if isNewKey {
			ts.size++
		}
		return isNewKey, nil
	}
	ts.size++
	node, err := newTrieNode(key[1:], value)
	if err != nil {
		return false, err
	}
	ts.links[next] = node
	return true, nil
}

//	getNode returns the node with the given key
func (ts *trieNode) getNode(key string) *trieNode {
	if len(key) == 0 {
		return ts
	}
	next := key[0]
	if nextNode, ok := ts.links[next]; ok {
		return nextNode.getNode(key[1:])
	}
	return nil
}

// contains returns true if the key is in the trie
// else returns false.
func (ts *trieNode) contains(key string) bool {
	return ts.getNode(key) != nil
}

// get returns the value for the key else
// returns nil, if the key is not in the trie.
func (ts *trieNode) get(key string) interface{} {
	if node := ts.getNode(key); node != nil {
		return node.value
	}
	return nil
}

// delete removed the value for the key
// and returns deleted as true if the key has a value,
// empty returns true if the entire path is deleted.
func (ts *trieNode) delete(key string) (deleted bool, empty bool) {
	deleted, empty = false, false
	if len(key) == 0 {
		deleted = ts.value != nil
		if deleted {
			ts.value = nil
			if ts.size == 0 {
				empty = true
			}
		}
	} else {
		next := key[0]
		if nextNode, ok := ts.links[next]; ok {
			deleted, empty = nextNode.delete(key[1:])
			if deleted {
				ts.size--
			}
			if empty {
				delete(ts.links, next)
				empty = ts.size == 0 && ts.value == nil
			}
		}
	}
	return
}

func (ts *trieNode) longestPrefixOf(s string, sIndex int) (result string) {
	result = ""
	if sIndex == len(s) {
		if ts.value != nil {
			result = s
		}
	} else {
		next := s[sIndex]
		if nextNode, ok := ts.links[next]; ok {
			result = nextNode.longestPrefixOf(s, sIndex+1)
			if result == "" && ts.value != nil {
				return s[:sIndex]
			}
		}
	}
	return
}

type nodeKeyValue struct {
	key   string
	value interface{}
}

type bytes []byte

func (a bytes) Len() int           { return len(a) }
func (a bytes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a bytes) Less(i, j int) bool { return a[i] < a[j] }

func (ts *trieNode) items(path []byte) <-chan nodeKeyValue {
	ch := make(chan nodeKeyValue, 1)
	go func() {
		if ts == nil {
			close(ch)
			return
		}
		if ts.value != nil {
			ch <- nodeKeyValue{key: string(path), value: ts.value}
		}
		var sortedKeys bytes
		sortedKeys = make([]byte, len(ts.links))
		for key := range ts.links {
			sortedKeys = append(sortedKeys, key)
		}
		sort.Sort(sortedKeys)
		for _, key := range sortedKeys {
			path = append(path, key)
			for kv := range ts.links[key].items(path) {
				ch <- kv
			}
			path = path[:len(path)-1]
		}
		close(ch)
	}()
	return ch
}

func (ts *trieNode) keys(path []byte) <-chan string {
	ch := make(chan string, 1)
	go func() {
		for keyValue := range ts.items(path) {
			ch <- keyValue.key
		}
		close(ch)
	}()
	return ch
}
