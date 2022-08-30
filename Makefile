build:
	go build -o ./build/goforbroke1006-proxy-tcp ./cmd/tcp
.PHONY: build

install:
	cp ./build/* ${GOPATH}/bin/
.PHONY: install