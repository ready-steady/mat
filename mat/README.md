# MATLAB MAT

The package provides an adapter to the
[MATLAB MAT-file API](http://www.mathworks.com/help/pdf_doc/matlab/apiext.pdf).

## Installation

Run:

```bash
$ go get github.com/go-math/format/mat
```

The above command will fail. However, it will properly clone this repository
into [`$GOPATH`](https://golang.org/doc/code.html#GOPATH). Go to that
directory:

```bash
$ cd $GOPATH/src/github.com/go-math/format/mat
```

Set the `MATLAB_ROOT` and `MATLAB_ARCH` environment variables according to your
MATLAB installation. It is preferable to do so in `~/.bash_profile` or alike.
For example:

```bash
export MATLAB_ROOT=/Applications/MATLAB_R2014a.app
export MATLAB_ARCH=maci64
```

Update the environment and run:

```bash
$ make install
```

In order to run an executable that relies on this package, the dynamic linker
should be able to discover MATLAB’s libraries. To this end, an appropriate
environment variable should be set depending on your system. For example, in
`.bash_profile`:

```bash
export DYLD_LIBRARY_PATH="$MATLAB_ROOT/bin/$MATLAB_ARCH:$DYLD_LIBRARY_PATH"
```
