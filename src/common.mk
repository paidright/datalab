GIT_SHA=$(shell git rev-parse HEAD)

.PHONY: version.go
version.go:
	echo "package main" > version.go
	echo "" >> version.go
	echo 'const currentVersion = "$(GIT_SHA)"' >> version.go
