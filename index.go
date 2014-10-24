package corpus

import (
	"encoding/json"

	"github.com/blevesearch/bleve"
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

// Index holds a Bleve index
type Index struct {
	index bleve.Index
}

func buildIndexMapping(language string) *bleve.IndexMapping {
	txtMapping := bleve.NewTextFieldMapping()
	txtMapping.Analyzer = language

	storeFieldOnlyMapping := bleve.NewTextFieldMapping()
	storeFieldOnlyMapping.Index = false
	storeFieldOnlyMapping.IncludeTermVectors = false
	storeFieldOnlyMapping.IncludeInAll = false

	docMapping := bleve.NewDocumentMapping()
	docMapping.AddSubDocumentMapping("id", bleve.NewDocumentDisabledMapping())
	docMapping.AddFieldMappingsAt("content", txtMapping)
	docMapping.AddFieldMappingsAt("title", txtMapping)
	docMapping.AddFieldMappingsAt("data", storeFieldOnlyMapping)

	mapping := bleve.NewIndexMapping()
	mapping.AddDocumentMapping("doc", docMapping)
	mapping.DefaultAnalyzer = language
	return mapping
}

func New(path, language string) (*Index, error) {
	index, err := bleve.Open(path)
	if err == bleve.ErrorIndexPathDoesNotExist {
		indexMapping := buildIndexMapping(language)
		index, err = bleve.New(path, indexMapping)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return &Index{index: index}, nil
}

type BleveDocument struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"title"`
	Data    string `json:"data"`
}

func (bd *BleveDocument) Type() string {
	return "doc"
}

func docToBleve(doc Document) (*BleveDocument, error) {
	data, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	return &BleveDocument{
		Id:      doc.Id(),
		Title:   doc.Title(),
		Content: doc.Content(),
		Data:    string(data),
	}, nil
}

func (b *Index) Insert(docs []Document) error {
	batch := bleve.NewBatch()
	for _, doc := range docs {
		bd, err := docToBleve(doc)
		if err != nil {
			return err
		}
		batch.Index(bd.Id, bd)
	}

	err := b.index.Batch(batch)
	return err
}

func (b *Index) Search(queryStr string, offset, limit int) (*bleve.SearchResult, error) {
	query := bleve.NewQueryStringQuery(queryStr)
	req := bleve.NewSearchRequestOptions(query, limit, offset, false)
	req.Highlight = bleve.NewHighlightWithStyle("ansi")
	return b.index.Search(req)
}

// Close releases the resources associated with the index.
func (b *Index) Close() {
	b.index.Close()
}
