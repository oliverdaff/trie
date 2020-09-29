package trie

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
