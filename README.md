# MATLAB MAT

The package provides an adapter to the [MATLAB MAT-file API][1].

## [Documentation][doc]

## Installation

Fetch the package:

```bash
go get -d github.com/ready-steady/mat
```

Go to the directory of the package:

```bash
cd $GOPATH/src/github.com/ready-steady/mat
```

Set the `MATLAB_ROOT` and `MATLAB_ARCH` environment variables according to your
MATLAB installation. For example,

```bash
export MATLAB_ROOT=/Applications/MATLAB_R2015a.app
export MATLAB_ARCH=maci64
```

Finally, install the package:

```bash
make install
```

In order to run an executable that relies on this package, the dynamic linker
should be able to discover MATLAB’s libraries. To this end, an appropriate
environment variable should be set depending on your system. For example,

```bash
export DYLD_LIBRARY_PATH="$MATLAB_ROOT/bin/$MATLAB_ARCH:$DYLD_LIBRARY_PATH"
```

## Contributing

1. Fork the project.
2. Implement your idea.
3. Create a pull request.

[1]: http://www.mathworks.com/help/pdf_doc/matlab/apiext.pdf

[doc]: http://godoc.org/github.com/ready-steady/mat
