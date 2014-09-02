all: $(patsubst %.tgffopt,%.tgff,$(wildcard *.tgffopt))

%.tgff: %.tgffopt
	tgff $*

clean:
	rm -f *.tgff *.eps *.vcg

.PHONY: all clean
