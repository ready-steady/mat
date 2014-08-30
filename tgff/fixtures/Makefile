inputs := $(wildcard *.tgffopt)
outputs := $(inputs:%.tgffopt=%.tgff)

all: $(outputs)

%.tgff: %.tgffopt
	tgff $*

clean:
	rm -f *.tgff *.eps *.vcg

.PHONY: all clean
