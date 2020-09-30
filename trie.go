package trie

import (
	"sort"
)

// trieNode is a internal representation of a trie.
//
type trieNode struct {
	links map[byte]*trieNode
	size  int
	value interface{}
}

func newTrieNode(key string, value interface{}, keyIndex int) *trieNode {
	if keyIndex == len(key) {
		return &trieNode{
			size:  0,
			value: value,
			links: make(map[byte]*trieNode),
		}
	}
	links := make(map[byte]*trieNode)
	links[key[keyIndex]] = newTrieNode(key, value, keyIndex+1)
	return &trieNode{
		size:  1,
		links: links,
	}
}

func (ts *trieNode) put(key string, value interface{}, keyIndex int) bool {
	if keyIndex == len(key) {
		isNewKey := ts.value == nil
		ts.value = value
		return isNewKey
	}
	next := key[keyIndex]
	if nextNode, ok := ts.links[next]; ok {
		isNewKey := nextNode.put(key, value, keyIndex+1)
		if isNewKey {
			ts.size++
		}
		return isNewKey
	}
	ts.size++
	ts.links[next] = newTrieNode(key, value, keyIndex+1)
	return true
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
