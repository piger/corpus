package corpus

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

type testDoc struct {
	DocId      string
	DocTitle   string
	DocContent string
}

func (d *testDoc) Id() string                   { return d.DocId }
func (d *testDoc) Title() string                { return d.DocTitle }
func (d *testDoc) Content() string              { return d.DocContent }
func (d *testDoc) MarshalJSON() ([]byte, error) { return json.Marshal(d.DocId) }

func Test_IndexAndSearch(t *testing.T) {
	tmpdir, _ := ioutil.TempDir("", "indexer_test_")
	defer os.RemoveAll(tmpdir)

	db := NewLucyDb(tmpdir, "en")
	defer db.Close()

	if err := db.Insert([]Document{
		&testDoc{"id1", "banana", "one"},
		&testDoc{"id2", "more banana", "two"},
		&testDoc{"id3", "", "three"},
	}); err != nil {
		t.Fatal(err)
	}

	nr, results := db.Search("three", 0, 100)
	if nr != 1 {
		t.Fatalf("Found %d results for 'three' instead of 1: %+v", nr, results)
	}

	nr, results = db.Search("banana", 0, 100)
	if nr != 2 {
		t.Fatalf("Found %d results for 'banana' instead of 1: %+v", nr, results)
	}
}
