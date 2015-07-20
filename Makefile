# Use libicu installation from homebrew
ifeq ($(shell uname), Darwin)
export CGO_LDFLAGS := -L/usr/local/opt/icu4c/lib 
export CGO_CFLAGS := -I/usr/local/opt/icu4c/include
endif

GO_FILES := *.go cmd/corpus/corpus.go
GO_FLAGS := -tags "libstemmer icu"

.PHONY: all install

all: install

install: $(GO_FILES)
	go install $(GO_FLAGS) ./cmd/corpus/

corpus: $(GO_FILES)
	go build $(GO_FLAGS) ./cmd/corpus/corpus.go
