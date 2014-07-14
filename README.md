corpus
======

Sample implementation of a small, self-contained search engine, built using [Apache Lucy](http://lucy.apache.org/).

An example command-line tool is provided, which will let you index and search a collection of text files.

Dependencies
------------

Required dependencies:

 * [Go](http://golang.org/)
 * *libmagic* (use ``apt-get install libmagic-dev`` to install it on a
   Debian system)
 * *Lucy* - we depend on a specific version of this C library, so a
   script is provided to install it (see below).

The default file indexing application uses some external binaries to
extract text from common file types. These dependencies are optional
and not required when using the code as a library:

 * *Lynx* to parse HTML
 * *pdftotext* for PDF files


Installation
------------

This software uses the
[golucy](https://github.com/philipsoutham/golucy) Go wrappers for
Lucy, which have only been tested successfully against a specific Lucy
version. As a result, installing from source is still a bit
complicated. To simplify the build process, you can use the
`install-lucy` script which automates this step.

### From sources

This is a step-by-step guide to install the software from the most
recent sources:

* Clone the source repository in the right place in your ``GOPATH``:
```
    $ mkdir -p $GOPATH/src/git.autistici.org/ale
    $ git clone https://git.autistici.org/ale/corpus.git \
          $GOPATH/src/git.autistici.org/ale/corpus
```

* Install Lucy using the provided script (this will ask you for your
  password as it's using ``sudo``):
```
    $ cd $GOPATH/src/git.autistici.org/ale/corpus
    $ ./install-lucy /usr/local/lucy
```

* Set the required environment variables so that Go can find the Lucy
  headers and shared libraries:
```
    $ CGO_LDFLAGS='-L/usr/local/lucy/lib -llucy -lcfish'
    $ CGO_CFLAGS=-I/usr/local/lucy/include
    $ LD_LIBRARY_PATH=/usr/local/lucy/lib
    $ export CGO_LDFLAGS CGO_CFLAGS LD_LIBRARY_PATH
```

* Finally build the ``corpus`` binary:
```
    $ go get ./...
    $ go build -v cmd/corpus/corpus.go
```

At the end of the above procedure, the `corpus` tool should have been
created in the current directory.

