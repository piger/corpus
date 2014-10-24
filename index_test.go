package corpus

import (
	"io/ioutil"
	"os"
	"path/filepath"
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

func Test_IndexAndSearch(t *testing.T) {
	tmpdir, _ := ioutil.TempDir("", "indexer_test_")
	defer os.RemoveAll(tmpdir)

	index, err := New(filepath.Join(tmpdir, "index"), "en")
	if err != nil {
		t.Fatalf("new Index error: %s\n", err)
	}
	defer index.Close()

	if err := index.Insert([]Document{
		&testDoc{"id1", "banana", "one"},
		&testDoc{"id2", "more banana", "two"},
		&testDoc{"id3", "", "three"},
	}); err != nil {
		t.Fatal(err)
	}

	results, err := index.Search("three", 0, 100)
	if err != nil {
		t.Fatalf("Search error: %s\n", err)
	} else if results.Total != 1 {
		t.Fatalf("Found %d results for 'three' instead of 1: %+v", results.Total, results)
	}
	resultdoc := results.Hits[0]
	if resultdoc.ID != "id3" {
		t.Errorf("bad result id: got=%s, want=id3", resultdoc.ID)
	}

	results, err = index.Search("banana", 0, 100)
	if err != nil {
		t.Fatalf("Search error: %s\n", err)
	} else if results.Total != 2 {
		t.Fatalf("Found %d results for 'banana' instead of 1: %+v", results.Total, results)
	}

	results, err = index.Search("boiler", 0, 100)
	if err != nil {
		t.Fatalf("Search error: %s\n", err)
	} else if results.Total != 0 {
		t.Fatalf("Found %d results for 'boiler' instead of none: %+v", results.Total, results)
	}
}
