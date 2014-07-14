package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestHtmlExtractor(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	h := `<html><body><h1>this is a <span class="big">test</h1></body></html>`
	s := "this is a test"
	fmt.Fprintf(f, h)
	f.Close()

	x := &htmlExtractor{}
	meta, err := x.Analyze(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(meta.Content) != s {
		t.Errorf("Content: %q (expected: %q)", meta.Content, s)
	}
}
