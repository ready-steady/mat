# MATLAB MAT-File

An adapter for the MATLAB MAT-file API. Refer to
[MATLAB External Interfaces](http://www.mathworks.com/help/pdf_doc/matlab/apiext.pdf)
for further details.

## Installation

Run:

```bash
$ go get github.com/gomath/format/mat
```

The above command will fail. However, it will properly clone this repository
into [`$GOPATH`](https://golang.org/doc/code.html#GOPATH). Go to that
directory:

```bash
$ cd $GOPATH/src/github.com/gomath/format/mat
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

## Contributing

1. [Fork](https://help.github.com/articles/fork-a-repo) the project.
2. Create a branch for your feature (`git checkout -b awesome-feature`).
3. Commit your changes (`git commit -am 'Add an awesome feature'`).
4. Push to the branch (`git push -u origin awesome-feature`).
5. [Create](https://help.github.com/articles/creating-a-pull-request)
   a pull request.
