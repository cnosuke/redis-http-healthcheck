NAME     := healthchecker
VERSION  := $(shell git describe --tags)
REVISION := $(shell git rev-parse --short HEAD)
SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

.PHONY: clean cross-build release-pack glide deps

bin/$(NAME): $(SRCS)
	go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME)

clean:
	rm -rf bin/* dist/* vendor/*

cross-build:
	for os in darwin linux; do \
		GOOS=$$os GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$$os-amd64/$(NAME); \
	done

release-pack: cross-build
	for os in darwin linux; do \
		zip -j dist/$$os-amd64.zip dist/$$os-amd64/$(NAME); \
	done

glide:
ifeq ($(shell command -v glide 2> /dev/null),)
	curl https://glide.sh/get | sh
endif

deps: glide
	glide install -v
