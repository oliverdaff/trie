# trie

[![PkgGoDev](https://pkg.go.dev/badge/github.com/oliverdaff/trie)](https://pkg.go.dev/github.com/oliverdaff/trie) [![Go Report Card](https://goreportcard.com/badge/github.com/oliverdaff/trie)](https://goreportcard.com/report/github.com/oliverdaff/trie) [![CircleCI](https://circleci.com/gh/oliverdaff/trie.svg?style=shield)](https://circleci.com/gh/oliverdaff/trie)

A trie written in Go.

A trie is a tree structure for storing strings with
common prefixes.

This implementation will store strings with values, the trie allow nil values to be stored but not nil.

# API

__Create new trie__

Creates a new empty trie.
```go
import trie
t := trie.NewTrie()
```


__Put key value__

Add a key and value in the trie.
```go
import trie
t := trie.NewTrie()
t.put("www.test.com", 1)
//Null value
var empty int
t.put("www.example.com", empty)
```

__Get value__

Get a value from the trie for a given key.

```go
import trie
t := trie.NewTrie()
t.put("www.test.com", 1)
t.get("www.test.com")
```

__Delete key__

Delete the key and value from the trie.

```go
import trie
t := trie.NewTrie()
t.Put("www.test.com", 1)
t.Delete("www.test.com")
```

__Contains key__

Check if the key exists in the trie.

```go
import trie
t := trie.NewTrie()
t.Put("www.test.com", 1)
t.Contains("www.test.com")
```

__Number of keys__

Get the number of keys in the trie.

```go
import trie
t := trie.NewTrie()
t.Put("www.test.com", 1)
t.Put("www.example.com", 1)
t.Size()
```

__No keys__

Find out if the trie has no keys.

```go
import trie
t := trie.NewTrie()
t.Put("www.test.com", 1)
t.IsEmpty()
```

__Find longest prefix__

Search the trie for the longest prefix of the given trie

```go
import trie
t := trie.NewTrie()
t.Put("www.test.example.com", 1)
t.Put("www.dev.example.com", 1)
t.LongestPrefixOf("www.test.example.com/home.html")
```

__Find keys with prefix__

Search the trie for keys with the given prefix.

```go
import trie
t := trie.NewTrie()
t.Put("www.test.example.com", 1)
t.Put("www.dev.example.com", 1)
t.KeysWithPrefix("www")
```

__Find all keys in the trie.__

Returns a channel with that return all the keys in the trie.

```go
import trie
t := trie.NewTrie()
t.Put("www.test.example.com", 1)
t.Put("www.dev.example.com", 1)
t.KeysWithPrefix("www")
```

__All key values in the trie__

__Find all keys in the trie.__

Returns a channel with that returns all the key values in the trie,
as `NodeKeyValue`.

```go
import trie
t := trie.NewTrie()
t.Put("www.test.example.com", 1)
t.Put("www.dev.example.com", 1)
t.Items()
```


## Tests
The tests can be invoked with `go test`

## License
MIT © Oliver Daff