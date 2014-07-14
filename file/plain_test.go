package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestPlainTextExtractor(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	s := "this is a test"
	fmt.Fprintf(f, s)
	f.Close()

	x := &plainTextExtractor{}
	meta, err := x.Analyze(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if meta.Content != s {
		t.Errorf("Content: %q (expected: %q)", meta.Content, s)
	}
}
