package corpus

import (
	"encoding/json"

	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index/upside_down"
	"github.com/rainycape/cld2"
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

func buildIndexMapping() *bleve.IndexMapping {
	mapping := bleve.NewIndexMapping()

	storeFieldOnlyMapping := bleve.NewTextFieldMapping()
	storeFieldOnlyMapping.Index = false
	storeFieldOnlyMapping.IncludeTermVectors = false
	storeFieldOnlyMapping.IncludeInAll = false

	itTextFieldMapping := bleve.NewTextFieldMapping()
	itTextFieldMapping.Analyzer = "my_it"

	enTextFieldMapping := bleve.NewTextFieldMapping()
	enTextFieldMapping.Analyzer = "my_en"

	genericTextFieldMapping := bleve.NewTextFieldMapping()
	genericTextFieldMapping.Analyzer = "my_base"

	docMapping := bleve.NewDocumentMapping()
	docMapping.DefaultAnalyzer = "my_base"
	docMapping.AddSubDocumentMapping("id", bleve.NewDocumentDisabledMapping())
	docMapping.AddFieldMappingsAt("content", genericTextFieldMapping)
	docMapping.AddFieldMappingsAt("title", genericTextFieldMapping)
	docMapping.AddFieldMappingsAt("data", storeFieldOnlyMapping)
	mapping.AddDocumentMapping("doc", docMapping)

	itDocMapping := bleve.NewDocumentMapping()
	itDocMapping.DefaultAnalyzer = "my_it"
	itDocMapping.AddSubDocumentMapping("id", bleve.NewDocumentDisabledMapping())
	itDocMapping.AddFieldMappingsAt("content", itTextFieldMapping)
	itDocMapping.AddFieldMappingsAt("title", genericTextFieldMapping)
	itDocMapping.AddFieldMappingsAt("data", storeFieldOnlyMapping)
	mapping.AddDocumentMapping("doc_it", itDocMapping)

	enDocMapping := bleve.NewDocumentMapping()
	enDocMapping.DefaultAnalyzer = "my_en"
	enDocMapping.AddSubDocumentMapping("id", bleve.NewDocumentDisabledMapping())
	enDocMapping.AddFieldMappingsAt("content", enTextFieldMapping)
	enDocMapping.AddFieldMappingsAt("title", genericTextFieldMapping)
	enDocMapping.AddFieldMappingsAt("data", storeFieldOnlyMapping)
	mapping.AddDocumentMapping("doc_en", enDocMapping)

	// mapping.DefaultAnalyzer = "my_en"
	mapping.DefaultField = "content"
	return mapping
}

func New(path string) (*Index, error) {
	index, err := bleve.Open(path)
	if err == bleve.ErrorIndexPathDoesNotExist {
		indexMapping := buildIndexMapping()
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
	Content string `json:"content"`
	Data    string `json:"data"`

	// this field was called "Type" but we need
	// to respect the bleve.Classifier interface!
	Kind string `json:"_type"`
}

func (bd *BleveDocument) Type() string {
	return bd.Kind
}

func docToBleve(doc Document) (*BleveDocument, error) {
	data, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	var docType string
	guessLang := cld2.Detect(doc.Content())
	fmt.Printf("LANG of %s: %s\n", doc.Title(), guessLang)
	if guessLang == "it" {
		docType = "doc_it"
	} else if guessLang == "en" {
		docType = "doc_en"
	} else {
		docType = "doc"
	}

	return &BleveDocument{
		Id:      doc.Id(),
		Title:   doc.Title(),
		Content: doc.Content(),
		Data:    string(data),
		Kind:    docType,
	}, nil
}

func (b *Index) Insert(docs []Document) error {
	batch := b.index.NewBatch()
	for _, doc := range docs {
		bd, err := docToBleve(doc)
		if err != nil {
			return err
		}
		if err := batch.Index(bd.Id, bd); err != nil {
			return err
		}
	}

	return b.index.Batch(batch)
}

func (b *Index) Search(queryStr string, offset, limit int, highlights bool) (*bleve.SearchResult, error) {
	query := bleve.NewQueryStringQuery(queryStr)
	req := bleve.NewSearchRequestOptions(query, limit, offset, false)
	if highlights {
		req.Highlight = bleve.NewHighlightWithStyle("ansi")
	}
	return b.index.Search(req)
}

// Close releases the resources associated with the index.
func (b *Index) Close() {
	b.index.Close()
}

func (b *Index) Dump() {
	for rowOrErr := range b.index.DumpAll() {
		switch rowOrErr := rowOrErr.(type) {
		case error:
			fmt.Printf("error dumping: %v\n", rowOrErr)
		case upside_down.UpsideDownCouchRow:
			fmt.Printf("%v\n", rowOrErr)
			fmt.Printf("Key:   % -100x\nValue: % -100x\n\n", rowOrErr.Key(), rowOrErr.Value())
		}
	}
}
