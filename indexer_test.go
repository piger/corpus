package corpus

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

type testDoc struct {
	DocId      string
	DocTitle   string
	DocContent string
}

func (d *testDoc) Id() string      { return d.DocId }
func (d *testDoc) Title() string   { return d.DocTitle }
func (d *testDoc) Content() string { return d.DocContent }
func (d *testDoc) ToJSON() string {
	data, _ := json.Marshal(d)
	return string(data)
}

func Test_IndexAndSearch(t *testing.T) {
	tmpdir, _ := ioutil.TempDir("", "indexer_test_")
	defer os.RemoveAll(tmpdir)

	db := NewLucyDb(tmpdir, "en")
	defer db.Close()

	if err := db.Insert([]Document{
		&testDoc{"id1", "title", "one"},
		&testDoc{"id2", "other title", "two"},
		&testDoc{"id3", "", "three"},
	}); err != nil {
		t.Fatal(err)
	}

	nr, results := db.Search("three", 0, 100)
	if nr != 1 {
		t.Fatalf("Found %d results instead of 1", nr)
	}
	log.Printf("%+v", results[0])
}
