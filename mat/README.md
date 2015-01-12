# MATLAB MAT

The package provides an adapter to the [MATLAB MAT-file API][1].

## [Documentation][doc]

## Installation

Run:

```bash
$ go get -d github.com/ready-steady/format/mat
```

Go to the directory of the package:

```bash
$ cd $GOPATH/src/github.com/ready-steady/format/mat
```

Set the `MATLAB_ROOT` and `MATLAB_ARCH` environment variables according to your
MATLAB installation. For example,

```bash
$ export MATLAB_ROOT=/Applications/MATLAB_R2014b.app
$ export MATLAB_ARCH=maci64
```

Run:

```bash
$ make install
```

In order to run an executable that relies on this package, the dynamic linker
should be able to discover MATLAB’s libraries. To this end, an appropriate
environment variable should be set depending on your system. For example,

```bash
$ export DYLD_LIBRARY_PATH="$MATLAB_ROOT/bin/$MATLAB_ARCH:$DYLD_LIBRARY_PATH"
```

[1]: http://www.mathworks.com/help/pdf_doc/matlab/apiext.pdf
[2]: https://golang.org/doc/code.html#GOPATH

[doc]: http://godoc.org/github.com/ready-steady/format/mat
