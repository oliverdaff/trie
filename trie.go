package trie

// trieNode is a internal representation of a trie.
//
type trieNode struct {
	links map[byte]*trieNode
	size  int
	value interface{}
}
