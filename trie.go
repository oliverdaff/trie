package trie

import (
	"errors"
	"sort"
)

// Trie can store strings
type Trie struct {
	root *trieNode
}

// NewTrie creates a new trie
func NewTrie() *Trie {
	node, _ := newTrieNode("", nil)
	return &Trie{root: node}
}

// Put inserts a string key into the trie
// with the given value.
// The key must not be empty and the value must
// not but nil.
func (ts *Trie) Put(key string, val interface{}) error {
	if key == "" {
		return errors.New("Can not put empty string as key")
	}
	_, err := ts.root.put(key, val)
	return err
}

// Get returns the value for the given key else nil.
func (ts *Trie) Get(key string) (interface{}, error) {
	if key == "" {
		return nil, errors.New("Can not get empty string as key")
	}
	return ts.root.get(key), nil
}

//Delete removes key if a value is found in the
//Trie associated with the key.
//Returns true if the key was deleted.
func (ts *Trie) Delete(key string) (bool, error) {
	if key == "" {
		return false, errors.New("Can not delete empty string as key")
	}
	deleted, _ := ts.root.delete(key)
	return deleted, nil
}

// Contains returns true if the key is in the trie.
func (ts *Trie) Contains(key string) (bool, error) {
	if key == "" {
		return false, errors.New("Can not delete empty string as key")
	}
	return ts.root.contains(key), nil
}

// Size returns the number of keys in the tri
func (ts *Trie) Size() int {
	return ts.root.size
}

// IsEmpty returns true if the Trie contains no keys
func (ts *Trie) IsEmpty() bool {
	return ts.Size() == 0
}

// LongestPrefixOf returns the longest key in the Trie that
// is a prefix of passed in key.
func (ts *Trie) LongestPrefixOf(key string) (string, error) {
	if key == "" {
		return "", errors.New("Can not search empty string as key")
	}
	return ts.root.longestPrefixOf(key, 0), nil

}

// KeysWithPrefix searches the Trie for all keys for which
// the prefix is a valid prefix.
func (ts *Trie) KeysWithPrefix(prefix string) <-chan string {
	prefixNode := ts.root.getNode(prefix)
	if prefixNode == nil {
		c := make(chan string, 0)
		close(c)
		return c
	}
	return prefixNode.keys(bytes(prefix))
}

// Keys returns all keys in the Trie
func (ts *Trie) Keys() <-chan string {
	return ts.root.keys(bytes(""))
}

// Items returns all the key value pairs in the Trie.
func (ts *Trie) Items() <-chan NodeKeyValue {
	return ts.root.items(make([]byte, 0))
}

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
	//create new node for remainder
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

// NodeKeyValue are the key and value pairs
// stored in the trie.
type NodeKeyValue struct {
	Key   string
	Value interface{}
}

type bytes []byte

func (a bytes) Len() int           { return len(a) }
func (a bytes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a bytes) Less(i, j int) bool { return a[i] < a[j] }

func (ts *trieNode) items(path []byte) <-chan NodeKeyValue {
	ch := make(chan NodeKeyValue, 1)
	go func() {
		if ts.value != nil {
			ch <- NodeKeyValue{Key: string(path), Value: ts.value}
		}
		var sortedKeys bytes
		sortedKeys = make([]byte, 0)
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
			ch <- keyValue.Key
		}
		close(ch)
	}()
	return ch
}
