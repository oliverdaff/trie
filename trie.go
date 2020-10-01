package trie

import (
	"sort"

	"github.com/pkg/errors"
)

// trieNode is a internal representation of a trie.
// Each node is root of its sub-trie. trieNode allows searching and adding new key-value pairs.
// Most operations, along with a string for the key, an index is passed, to mark the
// next character in the key that should be acted upon, rather than passing a substring with the first character removed.
// This is an optimization that allows keeping the asymptotic time required for each operation linear in the length of key.
type trieNode struct {
	links map[byte]*trieNode
	size  int
	value interface{}
}

// newTrieNode creates a new trie
//
// - key: the possibly empty key to store in the trie
// - value: the value to be associated with the key
// - keyIndex: the index of the start of the substring of the key to store in this subtrie.
func newTrieNode(key string, value interface{}, keyIndex int) (*trieNode, error) {
	if keyIndex > len(key) {
		return nil, errors.Errorf("Key index %d greater than key length %d for key %s", keyIndex, len(key), key)
	}
	if keyIndex < 0 {
		return nil, errors.Errorf("Key index less than 0 (%d) for key %s", keyIndex, key)
	}
	links := make(map[byte]*trieNode)
	if int(keyIndex) == len(key) {
		return &trieNode{
			size:  0,
			value: value,
			links: links,
		}, nil
	}
	node, err := newTrieNode(key, value, keyIndex+1)
	if err != nil {
		return nil, err
	}
	links[key[keyIndex]] = node

	return &trieNode{
		size:  1,
		links: links,
	}, nil

}

// put stores a key value pair in the trie.
// Returns true unless the key was already in the trie and got updated.
// A error is returned if the keyIndex is greater than the length of the key
// or the keyIndex is less than 0.
func (ts *trieNode) Put(key string, value interface{}) (bool, error) {
	return ts.put(key, value, 0)
}

func (ts *trieNode) put(key string, value interface{}, keyIndex int) (bool, error) {
	if keyIndex > len(key) {
		return false, errors.Errorf("Key index %d greater than key length %d for key %s", keyIndex, len(key), key)
	}
	if keyIndex < 0 {
		return false, errors.Errorf("Key index less than 0 (%d) for key %s", keyIndex, key)
	}
	if keyIndex == len(key) {
		isNewKey := ts.value == nil
		ts.value = value
		return isNewKey, nil
	}
	next := key[keyIndex]
	if nextNode, ok := ts.links[next]; ok {
		isNewKey, err := nextNode.put(key, value, keyIndex+1)
		if err != nil {
			return false, err
		}
		if isNewKey {
			ts.size++
		}
		return isNewKey, nil
	}
	ts.size++
	node, err := newTrieNode(key, value, keyIndex+1)
	if err != nil {
		return false, err
	}
	ts.links[next] = node
	return true, nil
}

func (ts *trieNode) getNode(key string, keyIndex int) *trieNode {
	if keyIndex == len(key) {
		return ts
	}
	next := key[keyIndex]
	if nextNode, ok := ts.links[next]; ok {
		return nextNode.getNode(key, keyIndex+1)
	}
	return nil
}

func (ts *trieNode) contains(key string) bool {
	return ts.getNode(key, 0) != nil
}

func (ts *trieNode) get(key string) interface{} {
	if node := ts.getNode(key, 0); node != nil {
		return node.value
	}
	return nil
}

func (ts *trieNode) delete(key string, keyIndex int) (deleted bool, empty bool) {
	deleted, empty = false, false
	if keyIndex == len(key) {
		deleted = ts.value == nil
		if deleted {
			ts.value = nil
			if ts.size == 0 {
				empty = true
			}
		}
	} else {
		next := key[keyIndex]
		if nextNode, ok := ts.links[next]; ok {
			deleted, empty = nextNode.delete(key, keyIndex+1)
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

func (ts *trieNode) longestPrefixOf(s string, sIndex int) (result *string) {
	result = nil
	if sIndex == len(s) {
		if ts.value != nil {
			result = &s
		}
	} else {
		next := s[sIndex]
		if nextNode, ok := ts.links[next]; ok {
			result = nextNode.longestPrefixOf(s, sIndex+1)
			if result == nil && !(ts.value == nil) {
				partial := s[:sIndex]
				result = &partial
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
