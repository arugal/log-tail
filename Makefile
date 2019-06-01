all: fmt build

build: log-tail

fmt:
	go fmt ./...

log-tail:
	go build -o bin/log-tail ./cmd

clean:
	rm -f ./bin/log-tail