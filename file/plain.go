package file

import "io/ioutil"

type plainTextExtractor struct{}

func (p *plainTextExtractor) Analyze(path string) (*Metadata, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// The document title is the relative path (minus the extension).
	return &Metadata{
		Title:   TitleFromPath(path),
		Content: string(data),
	}, nil
}

func (p *plainTextExtractor) String() string { return "text" }

func init() {
	RegisterExtractor("text/plain", &plainTextExtractor{})

	// Other, more specific text/something extractors will take
	// precendence over this one.
	RegisterExtractor("text/", &plainTextExtractor{})
}
