package corpus

import (
	"github.com/philipsoutham/golucy/v0.0.1"
)

func NewStoredField(name string) *golucy.Field {
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

type LucyDb struct {
	schema *golucy.Schema
	index  *golucy.Index
}

func NewLucyDb(path, language string) *LucyDb {
	schema := golucy.NewSchema()
	schema.AddField(golucy.NewIdField("id"))
	schema.AddField(golucy.NewFTField("content", language))
	schema.AddField(golucy.NewFTField("title", language))
	schema.AddField(NewStoredField("data"))

	index := golucy.NewIndex(path, true, false, schema)

	return &LucyDb{schema, index}
}

func docToLucy(doc Document) golucy.Document {
	return golucy.Document{
		"id":      doc.Id(),
		"title":   doc.Title(),
		"content": doc.Content(),
		"data":    doc.ToJSON(),
	}
}

func (l *LucyDb) Insert(docs []Document) error {
	writer := l.index.NewIndexWriter()
	defer writer.Close()

	lucyDocs := make([]golucy.Document, 0, len(docs))
	for _, doc := range docs {
		lucyDocs = append(lucyDocs, docToLucy(doc))
	}

	writer.AddDocs(lucyDocs...)
	writer.Commit()

	return nil
}

func (l *LucyDb) Search(queryStr string, offset, limit uint) (uint, []string) {
	reader := l.index.NewIndexReader()
	defer reader.Close()

	query := reader.ParseQuery(queryStr)
	defer query.Close()

	// Run the query but only return the full object data.
	numResults, results := reader.Search(query, offset, limit, "id", "data", false)
	out := make([]string, 0, len(results))
	for _, r := range results {
		out = append(out, r.Text)
	}
	return numResults, out
}

func (l *LucyDb) Close() {
	l.index.Close()
	l.schema.Close()
}
