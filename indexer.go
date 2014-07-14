package corpus

import (
	"git.autistici.org/ale/corpus/third_party/golucy"
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

type Db struct {
	schema *golucy.Schema
	index  *golucy.Index
}

func NewLucyDb(path, language string) *Db {
	schema := golucy.NewSchema()
	schema.AddField(golucy.NewIdField("id"))
	schema.AddField(golucy.NewFTField("content", language, true))
	schema.AddField(golucy.NewFTField("title", language, true))
	schema.AddField(NewStoredField("data"))

	index := golucy.NewIndex(path, true, false, schema)

	return &Db{schema, index}
}

func docToLucy(doc Document) (golucy.Document, error) {
	data, err := doc.MarshalJSON()
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

func (l *Db) Insert(docs []Document) error {
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

func (l *Db) Search(queryStr string, offset, limit uint) (uint, []string) {
	reader := l.index.NewIndexReader()
	defer reader.Close()

	query := reader.ParseQuery(queryStr, true)
	defer query.Close()

	// Run the query but only return the full object data.
	numResults, results := reader.Search(query, offset, limit, "id", "data", false)
	out := make([]string, 0, len(results))
	for _, r := range results {
		out = append(out, r.Text)
	}
	return numResults, out
}

func (l *Db) Close() {
	l.index.Close()
	l.schema.Close()
}
