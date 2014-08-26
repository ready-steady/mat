ifndef MATLAB_ROOT
$(error MATLAB_ROOT should be defined.)
endif

ifndef MATLAB_ARCH
$(error MATLAB_ARCH should be defined.)
endif

CGO_LDFLAGS := -L$(MATLAB_ROOT)/bin/$(MATLAB_ARCH)
CGO_CFLAGS := -I$(MATLAB_ROOT)/extern/include

GO_FLAGS += CGO_CFLAGS='$(CGO_CFLAGS)'
GO_FLAGS += CGO_LDFLAGS='$(CGO_LDFLAGS)'

all: build

build test install:
	$(GO_FLAGS) go $@

.PHONY: all build test install
