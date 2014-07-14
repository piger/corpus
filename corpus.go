package corpus

// A Document is the basic representation of an indexable object.
// It can have a "title", and some "content", which will be used
// for indexing with our simple tiny schema.
//
// Documents must be JSON-serializable, i.e. either a base type, or an
// implementation of the json.Marshaler interface (even though this is
// not explicitly set in this interface due to the "base type" case).
type Document interface {
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
