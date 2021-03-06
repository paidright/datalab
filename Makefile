name = datalab

.PHONY: lint vet test check clean
TOOLS = `go list ./src/... | grep -v /vendor/`

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint:
	golint $(TOOLS)

# http://godoc.org/code.google.com/p/go.tools/cmd/vet
# go get code.google.com/p/go.tools/cmd/vet
vet:
	go vet $(TOOLS)

check: test vet lint

test:
	go test -cover $(TOOLS)

clean:
	-rm -r dist/github.com/paidright/datalab/*

dist/github.com/paidright/datalab/sort_csv: ./sort_csv.sh ./sort_csv_asc.sh ./sort_csv_desc.sh
	-mkdir -p ./dist/github.com/paidright/datalab
	cp sort_csv.sh dist/github.com/paidright/datalab/sort_csv
	cp sort_csv_asc.sh dist/github.com/paidright/datalab/sort_csv_asc
	cp sort_csv_desc.sh dist/github.com/paidright/datalab/sort_csv_desc

dist: clean dist/github.com/paidright/datalab/sort_csv
	echo $(TOOLS)
	set -e; \
	for dir in $(TOOLS); do \
		ls; \
		cd $$GOPATH/src/$$dir && make version.go && cd ../..; \
		echo building $$dir; \
		packr2 build -o ./dist/`echo $$dir | sed s/src//` $$dir; \
		GOOS=darwin GOARCH=amd64 packr2 build -o ./dist/`echo $$dir | sed s/src//`-mac $$dir; \
	done

docker_image_build: build
