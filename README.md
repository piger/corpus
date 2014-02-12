corpus
======

Sample implementation of a small, self-contained search engine, built
using [Apache Lucy](http://lucy.apache.org/).

An example command-line tool is provided, which will let you index
and search a collection of text files.

Building
--------

This software uses the
[golucy](https://github.com/philipsoutham/golucy) Go wrappers for
Lucy, which have only been tested successfully against a specific Lucy
version. To simplify the build process, we've provided the
`install-lucy` script. Use it this way (once you have downloaded the
`corpus` code in the appropriate location in your GOPATH):

    $ ./install-lucy /usr/local/lucy
    $ CGO_LDFLAGS='-L/usr/local/lucy/lib -llucy -lcfish'
    $ CGO_CFLAGS=-I/usr/local/lucy/include
    $ LD_LIBRARY_PATH=/usr/local/lucy/lib
    $ export CGO_LDFLAGS CGO_CFLAGS LD_LIBRARY_PATH
    $ go get github.com/philipsoutham/golucy
    $ go build -v cmd/corpus/corpus.go

This should create the `corpus` tool in the current directory.