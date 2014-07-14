package corpus

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
)

func createTestFs(fs []string) string {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}

	for _, filename := range fs {
		p := path.Join(dir, filename)
		os.MkdirAll(path.Dir(p), 0700)
		ioutil.WriteFile(p, []byte("test"), 0600)
	}

	return dir
}

func doTestWalker(t *testing.T, w *Walker, fs []string, expected []string) {
	dir := createTestFs(fs)
	defer os.RemoveAll(dir)

	var result []string
	err := w.Walk(dir, func(p string, info os.FileInfo, err error) error {
		result = append(result, strings.TrimPrefix(p, dir)[1:])
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Bad result from Walk(): got=%v, want=%v", result, expected)
	}
}

func TestWalker_Base(t *testing.T) {
	var w Walker
	doTestWalker(t, &w,
		[]string{"dir1/testfile", "dir2/testfile"},
		[]string{"dir1/testfile", "dir2/testfile"})
}

func TestWalker_ExcludeDirectory(t *testing.T) {
	w := Walker{
		Exclude: []string{"dir1"},
	}
	doTestWalker(t, &w,
		[]string{"dir1/testfile", "dir2/testfile"},
		[]string{"dir2/testfile"})
}

func TestWalker_Exclude(t *testing.T) {
	w := Walker{
		Exclude: []string{"*file1"},
	}
	doTestWalker(t, &w,
		[]string{"testfile1", "testfile2", "dir1/testfile1"},
		[]string{"testfile2"})
}

func TestWalker_Include(t *testing.T) {
	w := Walker{
		Include: []string{"*.txt"},
	}
	doTestWalker(t, &w,
		[]string{"dir1/subdir/test.txt", "dir2.txt/testfile"},
		[]string{"dir1/subdir/test.txt"})
}

func TestWalker_IncludeExclude(t *testing.T) {
	w := Walker{
		Exclude: []string{"dir2"},
		Include: []string{"*.txt"},
	}
	doTestWalker(t, &w,
		[]string{"dir1/test.txt", "dir2/test.txt"},
		[]string{"dir1/test.txt"})
}
