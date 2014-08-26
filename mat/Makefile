# In order to compile, MATLAB_ROOT and MATLAB_ARCH are required to be set. For
# example:
#
# export MATLAB_ROOT=/Applications/MATLAB_R2014a.app
# export MATLAB_ARCH=maci64
#
# In order to run, the MATLAB libraries are required to be accessible for the
# dynamic linker. For example:
#
# export DYLD_LIBRARY_PATH="$MATLAB_ROOT/bin/$MATLAB_ARCH:$DYLD_LIBRARY_PATH"

ifndef MATLAB_ROOT
$(error MATLAB_ROOT should be undefined.)
endif

ifndef MATLAB_ARCH
$(error MATLAB_ARCH should be undefined.)
endif

CGO_LDFLAGS := -L$(MATLAB_ROOT)/bin/$(MATLAB_ARCH)
CGO_CFLAGS := -I$(MATLAB_ROOT)/extern/include

GO_FLAGS += CGO_CFLAGS='$(CGO_CFLAGS)'
GO_FLAGS += CGO_LDFLAGS='$(CGO_LDFLAGS)'

all: build

build test:
	$(GO_FLAGS) go $@

.PHONY: all build test
