tzcon: *.go
	go build

install:
	cp tzcon ~/bin/tzcon

.PHONY: \
	install
