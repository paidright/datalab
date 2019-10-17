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

$(PLATFORMS): dist/marx_darwin_finder_build
	@echo "Building for $(os)-$(arch)"
	@-mkdir -p dist/$(os)-$(arch)
	@-rm -r dist/$(os)-$(arch)
	GOOS=$(os) GOARCH=$(arch) go build -o 'dist/marx_$(os)_$(arch)' .
	@chmod +x dist/marx_$(os)_$(arch)
	@if [ $(os) = windows ]; then mv dist/marx_$(os)_$(arch) dist/marx_$(os)_$(arch).exe; fi

dist/marx_darwin_finder_build:
	GOOS=darwin GOARCH=amd64 go build --tags munge_cwd -o 'dist/marx_darwin_finder_build' .

test:
	bash .envrc && go test
