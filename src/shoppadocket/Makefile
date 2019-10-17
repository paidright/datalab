include ../common.mk

PLATFORMS := linux/amd64 windows/amd64 linux/arm darwin/amd64
temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
GIT_SHA=$(shell git rev-parse HEAD)
SOURCES==$(wildcard *.go)

release: version.go
	make -l inner_release

.PHONY: inner_release
inner_release: $(PLATFORMS)

$(PLATFORMS): .got_packr_bin
	@echo "Building for $(os)-$(arch)"
	@-mkdir -p dist/$(os)-$(arch)
	@-rm -r dist/$(os)-$(arch)
	GOOS=$(os) GOARCH=$(arch) packr2 build -o 'dist/shoppadocket_$(os)_$(arch)' .
	@chmod +x dist/shoppadocket_$(os)_$(arch)
	@if [ $(os) = windows ]; then mv dist/shoppadocket_$(os)_$(arch) dist/shoppadocket_$(os)_$(arch).exe; fi

.got_packr_bin:
	go get -u github.com/gobuffalo/packr/v2/packr2
	touch .got_packr_bin
