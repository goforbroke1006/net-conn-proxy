all: build
.PHONY: all

build:
	go build ./
.PHONY: build

install:
	@bash ./bin/install.sh
.PHONY: install
