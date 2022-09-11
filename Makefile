all: build
.PHONY: all

build:
	go build ./
.PHONY: build

install:
	go install ./
.PHONY: install