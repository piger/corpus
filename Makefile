GO_FILES := *.go cmd/corpus/corpus.go
GO_FLAGS := -tags "libstemmer"

.PHONY: all install

all: install

install: $(GO_FILES)
	go install $(GO_FLAGS) ./cmd/corpus/

corpus: $(GO_FILES)
	go build $(GO_FLAGS) ./cmd/corpus/corpus.go
