package file

import "os/exec"

type pdfExtractor struct{}

func (p *pdfExtractor) Analyze(path string) (*Metadata, error) {
	data, err := exec.Command("pdftotext", "-q", "-nopgbrk", "-enc", "UTF-8", path).Output()
	if err != nil {
		return nil, err
	}
	return &Metadata{
		Title:   TitleFromPath(path),
		Content: string(data),
	}, nil
}

func (p *pdfExtractor) String() string { return "pdf" }

func init() {
	RegisterExtractor("application/pdf", &pdfExtractor{})
	RegisterExtractor("text/pdf", &pdfExtractor{})
}
