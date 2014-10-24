corpus
======

Sample implementation of a small, self-contained search engine.

An example command-line tool is provided, which will let you index and search a collection of text files.

Dependencies
------------

Required dependencies:

 * [Go](http://golang.org/)
 * *libmagic* (use ``apt-get install libmagic-dev`` to install it on a
   Debian system)
 * *libicu* (use ``apt-get install libicu-dev`` to install it on a Debian
   system)

To install the dependencies on OS X with homebrew: ``brew install libmagic icu4c``.

The default file indexing application uses some external binaries to
extract text from common file types. These dependencies are optional
and not required when using the code as a library:

 * *Lynx* to parse HTML
 * *pdftotext* for PDF files


Installation
------------

This is a step-by-step guide to install the software from the most
recent sources:

```
$ go get -tags "libstemmer icu" git.autistici.org/ale/corpus/cmd/corpus
```

At the end of the above procedure, the `corpus` tool should have been
created inside ``$GOPATH/bin/``.

A GNU makefile is provided for compiling on OS X.

