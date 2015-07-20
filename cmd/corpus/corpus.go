package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/piger/corpus"
	"github.com/piger/corpus/file"
)

var (
	dbPath     = flag.String("db", "./db", "Path to the index")
	doSearch   = flag.Bool("search", true, "Search for something (default)")
	doIndex    = flag.Bool("index", false, "Index documents")
	doDump     = flag.Bool("dump", false, "Dump index database")
	limit      = flag.Int("limit", 20, "Limit number of search results")
	noScores   = flag.Bool("no-score", false, "Hide score in results")
	highlights = flag.Bool("highlights", false, "Show highlights from search results")

	// Multi-valued flags.
	includes strslice
	excludes = strslice{
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

func runSearch(index *corpus.Index, args []string) {
	// Execute query and display results.
	results, err := index.Search(strings.Join(args, " "), 0, *limit, *highlights)
	if err != nil {
		log.Fatalf("Search error: %s\n", err)
	} else if results.Total == 0 {
		log.Printf("No results")
		return
	}
	for _, hit := range results.Hits {
		if *noScores {
			fmt.Printf("%s\n", hit.ID)
		} else {
			fmt.Printf(" %-6.4f  %s\n", hit.Score, hit.ID)
		}
		if *highlights {
			hl := ""
			for _, fragments := range hit.Fragments {
				for _, fragment := range fragments {
					hl += fmt.Sprintf("%s", fragment)
				}
			}
			fmt.Printf("%s\n", hl)
		}
	}
}

func runIndex(index *corpus.Index, args []string) {
	docs := make([]corpus.Document, 0)

	w := &corpus.Walker{
		Exclude: excludes,
		Include: includes,
		MinSize: 1024,
	}

	// For each argument, process it or recurse if it's a directory.
	for _, root := range args {
		log.Printf("scanning %s ...", root)
		err := w.Walk(root, func(path string, info os.FileInfo, fileErr error) error {
			if fileErr == nil {
				doc, err := file.New(path)
				if err != nil {
					log.Printf("%s: %v", path, err)
				} else {
					docs = append(docs, doc)
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("Cannot scan %s: %s", root, err)
		}
	}

	if len(docs) == 0 {
		log.Fatal("No documents found!")
	}

	// Now add all documents to the index.
	if err := index.Insert(docs); err != nil {
		log.Fatal(err)
	}

	log.Printf("added %d documents", len(docs))
}

func main() {
	flag.Var(&excludes, "exclude", "Exclude pattern")
	flag.Var(&includes, "include", "Include pattern")
	flag.Parse()

	if !*doSearch && !*doIndex {
		log.Fatal("You have to specify one of --search or --index")
	}
	if flag.NArg() < 1 {
		log.Fatal("Not enough arguments")
	}

	index, err := corpus.New(*dbPath)
	if err != nil {
		log.Fatalf("Cannot open index directory: %s\n", err)
	}
	defer index.Close()

	if *doIndex {
		runIndex(index, flag.Args())
	} else if *doDump {
		index.Dump()
	} else {
		runSearch(index, flag.Args())
	}
}
