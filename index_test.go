package corpus

import (
	"io/ioutil"
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

func Test_IndexAndSearch(t *testing.T) {
	tmpdir, _ := ioutil.TempDir("", "indexer_test_")
	defer os.RemoveAll(tmpdir)

	index := New(tmpdir, "en")
	defer index.Close()

	if err := index.Insert([]Document{
		&testDoc{"id1", "banana", "one"},
		&testDoc{"id2", "more banana", "two"},
		&testDoc{"id3", "", "three"},
	}); err != nil {
		t.Fatal(err)
	}

	nr, results := index.Search("three", 0, 100)
	if nr != 1 {
		t.Fatalf("Found %d results for 'three' instead of 1: %+v", nr, results)
	}
	if !results.Next() {
		t.Fatal("Result.Next() failed")
	}
	var resultdoc testDoc
	if err := results.Value(&resultdoc); err != nil {
		t.Fatal("Result.Value(): ", err)
	}
	if resultdoc.DocId != "id3" {
		t.Errorf("bad result id: got=%s, want=id3", resultdoc.DocId)
	}

	nr, results = index.Search("banana", 0, 100)
	if nr != 2 {
		t.Fatalf("Found %d results for 'banana' instead of 1: %+v", nr, results)
	}

	nr, results = index.Search("boiler", 0, 100)
	if nr != 0 {
		t.Fatalf("Found %d results for 'boiler' instead of none: %+v", nr, results)
	}
}
