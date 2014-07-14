package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"git.autistici.org/ale/corpus"
)

var (
	dbPath   = flag.String("db", "./db", "Path to the index")
	doSearch = flag.Bool("search", true, "Search for something (default)")
	doIndex  = flag.Bool("index", false, "Index documents")
	limit    = flag.Int("limit", 20, "Limit number of search results")
	lang     = flag.String("lang", "en", "Language (for indexing)")

	// Multi-valued flags.
	includes []string
	excludes = []string{
		".*", "*~", "*.bak",
	}
)

type strslice []string

func (s *strslice) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *strslice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func init() {
	flag.Var(&excludes, "exclude", "Exclude pattern")
	flag.Var(&includes, "include", "Include pattern")
}

type FsDoc struct {
	Path string
}

func (d *FsDoc) Id() string    { return d.Path }
func (d *FsDoc) Title() string { return "" }

func (d *FsDoc) Content() string {
	data, err := ioutil.ReadFile(d.Path)
	if err != nil {
		return ""
	}
	return string(data)
}

func (d *FsDoc) ToJSON() string {
	data, _ := json.Marshal(d)
	return string(data)
}

func DecodeDoc(data string) *FsDoc {
	var d FsDoc
	json.Unmarshal([]byte(data), &d)
	return &d
}

func search(db corpus.Index, args []string) {
	// Execute query and display results.
	_, results := db.Search(strings.Join(args, " "), 0, uint(*limit))
	for _, r := range results {
		doc := DecodeDoc(r)
		fmt.Println(doc.Path)
	}
}

func index(db corpus.Index, args []string) {
	docs := make([]corpus.Document, 0)

	w := &corpus.Walker{
		Exclude: excludes,
		Include: includes,
		MinSize: 1024,
	}

	// For each argument, process it or recurse if it's a directory.
	for _, root := range args {
		err := w.Walk(root, func(path string, info os.FileInfo, err error) error {
			docs = append(docs, &FsDoc{path})
			return nil
		})
		if err != nil {
			log.Printf("Cannot scan %s: %s", root, err)
		}
	}

	if len(docs) == 0 {
		log.Fatal("No documents found!")
	}

	// Now add all documents to the db.
	if err := db.Insert(docs); err != nil {
		log.Fatal(err)
	}

	log.Printf("added %d documents", len(docs))
}

func main() {
	flag.Parse()

	if !*doSearch && !*doIndex {
		log.Fatal("You have to specify one of --search or --index")
	}
	if flag.NArg() < 1 {
		log.Fatal("Not enough arguments")
	}

	if fi, err := os.Stat(*dbPath); err != nil {
		if err := os.MkdirAll(*dbPath, 0755); err != nil {
			log.Fatal(err)
		}
	} else if !fi.IsDir() {
		log.Fatal("Index directory already exists and is not a directory")
	}

	db := corpus.NewLucyDb(*dbPath, *lang)
	defer db.Close()

	if *doIndex {
		index(db, flag.Args())
	} else {
		search(db, flag.Args())
	}
}
