package corpus

import (
	"errors"
	"encoding/json"
	"io/ioutil"
	"mime"
	"path/filepath"
	"strings"
)

var ErrUnknownMimeType = errors.New("unknown/unsupported mime type")

type Extractor interface {
	Analyze(path string) (Document, error)
}

var extractors = make(map[string]Extractor)

func RegisterExtractor(mimetype string, e Extractor) {
	extractors[mimetype] = e
}

func primaryType(mimetype string) string {
	if strings.Contains(mimetype, ";") {
		return strings.SplitN(mimetype, ";", 2)[0]
	}
	return mimetype
}

func Analyze(path string) (Document, error) {
	mimetype := primaryType(mime.TypeByExtension(filepath.Ext(path)))
	ext, ok := extractors[mimetype]
	if !ok {
		return nil, ErrUnknownMimeType
	}
	return ext.Analyze(path)
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

func (d *FsDoc) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Path)
}

func (d *FsDoc) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &d.Path)
}

type plainTextExtractor struct{}

func (p *plainTextExtractor) Analyze(path string) (Document, error) {
	return &FsDoc{path}, nil
}
