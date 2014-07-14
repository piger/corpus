package file

import "os/exec"

type htmlExtractor struct{}

func (h *htmlExtractor) Analyze(path string) (*Metadata, error) {
	data, err := exec.Command("lynx", "-dump", "-force_html", "-nolist", path).Output()
	if err != nil {
		return nil, err
	}
	return &Metadata{
		Title:   TitleFromPath(path),
		Content: string(data),
	}, nil
}

func (h *htmlExtractor) String() string { return "html" }

func init() {
	RegisterExtractor("text/html", &htmlExtractor{})
}
