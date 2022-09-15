all: build
.PHONY: all

build:
	go build ./
.PHONY: build

install:
	@bash ./bin/install.sh
.PHONY: install

release:
	GOOS=linux   GOARCH=amd64 go build -o ./release/net-conn-proxy--linux-x64       ./
	GOOS=linux   GOARCH=386   go build -o ./release/net-conn-proxy--linux-x86       ./
	GOOS=darwin  GOARCH=amd64 go build -o ./release/net-conn-proxy--darwin-x64      ./
	GOOS=windows GOARCH=amd64 go build -o ./release/net-conn-proxy--windows-x64.exe ./
	GOOS=windows GOARCH=386   go build -o ./release/net-conn-proxy--windows-x86.exe ./
