package file

import (
	"path/filepath"
	"strings"
)

type Metadata struct {
	Title   string
	Content string `json:"-"`
}

type Document struct {
	Path string
	Meta *Metadata
}

func (d *Document) Id() string      { return d.Path }
func (d *Document) Title() string   { return d.Meta.Title }
func (d *Document) Content() string { return d.Meta.Content }

func New(path string) (*Document, error) {
	meta, err := Analyze(path)
	if err != nil {
		return nil, err
	}
	return &Document{Path: path, Meta: meta}, nil
}

func TitleFromPath(path string) string {
	return strings.Replace(strings.TrimSuffix(path, filepath.Ext(path)), "/", " ", -1)
}
