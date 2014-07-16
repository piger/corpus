package corpus

import (
	"encoding/json"

	"git.autistici.org/ale/corpus/third_party/golucy"
)

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

// Index holds a Lucy index and an associated schema.
type Index struct {
	schema *golucy.Schema
	index  *golucy.Index
}

// New returns a new index. Indexes are mono-lingual, and their
// language must be specified at creation time.
func New(path, language string) *Index {
	schema := golucy.NewSchema()
	schema.AddField(golucy.NewIdField("id"))
	schema.AddField(golucy.NewFTField("content", language, true))
	schema.AddField(golucy.NewFTField("title", language, true))
	schema.AddField(newStoredField("data"))

	index := golucy.NewIndex(path, true, false, schema)

	return &Index{schema, index}
}

func docToLucy(doc Document) (golucy.Document, error) {
	data, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}
	return golucy.Document{
		"id":      doc.Id(),
		"title":   doc.Title(),
		"content": doc.Content(),
		"data":    string(data),
	}, nil
}

// Insert adds some documents to the index as a single batch
// operation.
func (l *Index) Insert(docs []Document) error {
	writer := l.index.NewIndexWriter()
	defer writer.Close()

	lucyDocs := make([]golucy.Document, 0, len(docs))
	for _, doc := range docs {
		ldoc, err := docToLucy(doc)
		if err != nil {
			return err
		}
		lucyDocs = append(lucyDocs, ldoc)
	}

	writer.AddDocs(lucyDocs...)
	writer.Commit()

	return nil
}

// An iterator scans through search results, providing a way to
// deserialize the associated object data. Use it like this:
//
//     _, iter := index.Search("query")
//     for iter.Next() {
//             var doc MyType
//             iter.Value(&doc)
//     }
//
type Iterator struct {
	pos        int
	numResults int
	results    []*golucy.SearchResult
}

func (i *Iterator) Next() bool {
	i.pos++
	return i.pos < i.numResults
}

func (i *Iterator) Value(obj interface{}) error {
	return json.Unmarshal([]byte(i.results[i.pos].Text), obj)
}

func (i *Iterator) Score() float32 {
	return i.results[i.pos].Score
}

func (i *Iterator) MatchedTerms() []string {
	return i.results[i.pos].MatchedTerms
}

func (l *Index) Search(queryStr string, offset, limit uint) (int, *Iterator) {
	reader := l.index.NewIndexReader()
	defer reader.Close()

	query := reader.ParseQuery(queryStr, true)
	defer query.Close()

	// Run the query but only return the full object data.
	numResults, results := reader.Search(query, offset, limit, "id", "data", true)
	return int(numResults), &Iterator{-1, int(numResults), results}
}

// Close releases the resources associated with the index.
func (l *Index) Close() {
	l.index.Close()
	l.schema.Close()
}

// Create a new Lucy field which is stored but not indexed (for the
// opaque document representation).
func newStoredField(name string) *golucy.Field {
	return &golucy.Field{
		Name:      name,
		IndexType: golucy.StringType,
		IndexOptions: &golucy.IndexOptions{
			Indexed:       false,
			Stored:        true,
			Sortable:      false,
			Highlightable: false,
		},
	}
}
