package file

import (
	"errors"
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/rakyll/magicmime"
)

var ErrUnknownMimeType = errors.New("unknown/unsupported mime type")

type Extractor interface {
	Analyze(path string) (*Metadata, error)
}

// The default extractor for unknown file types simply stores the path
// as the title.
type defaultExtractor struct{}

func (d *defaultExtractor) Analyze(path string) (*Metadata, error) {
	return &Metadata{
		Title:   TitleFromPath(path),
		Content: "",
	}, nil
}

func (d *defaultExtractor) String() string { return "default" }

// Registry of known extractors.
type extractorEntry struct {
	mimetype  string
	extractor Extractor
}

var (
	extractors []extractorEntry
	once       sync.Once
)

// Sort the extractor list in order of decreasing length of the
// 'mimetype' field (for prefix matches).
type extractorList []extractorEntry

func (l extractorList) Len() int      { return len(l) }
func (l extractorList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l extractorList) Less(i, j int) bool {
	return len(l[i].mimetype) > len(l[j].mimetype)
}

func sortExtractors() {
	sort.Sort(extractorList(extractors))
}

func RegisterExtractor(mimetype string, e Extractor) {
	extractors = append(extractors, extractorEntry{mimetype, e})
}

func primaryType(mimetype string) string {
	if strings.Contains(mimetype, ";") {
		return strings.SplitN(mimetype, ";", 2)[0]
	}
	return mimetype
}

var mm *magicmime.Magic

func init() {
	var err error
	mm, err = magicmime.New(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_ERROR | magicmime.MAGIC_SYMLINK)
	if err != nil {
		panic(err)
	}

	RegisterExtractor("", &defaultExtractor{})
}

func Analyze(path string) (*Metadata, error) {
	mimetype, err := mm.TypeByFile(path)
	if err != nil {
		return nil, err
	}
	mimetype = primaryType(mimetype)
	once.Do(sortExtractors)
	for _, t := range extractors {
		if strings.HasPrefix(mimetype, t.mimetype) {
			log.Printf("%s -> %s (%v)", path, mimetype, t.extractor)
			return t.extractor.Analyze(path)
		}
	}
	return nil, ErrUnknownMimeType
}
