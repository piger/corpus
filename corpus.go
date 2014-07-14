package corpus

import "encoding/json"

// A Document is the basic representation of an indexable object.
// It can have a "title", and some "content", which will be used
// for indexing with our simple tiny schema.
type Document interface {
	json.Marshaler

	Id() string
	Title() string
	Content() string
}

// High-level interface to the index.
type Index interface {
	Close()
	Insert(docs []Document) error
	Search(queryStr string, offset, limit uint) (uint, []string)
}
