export PATH := $(GOPATH)/bin:$(PATH)

all: clean fmt build

build: log-tail

fmt:
	go fmt ./...

file:
	rm -rf ./assets/tail/static/*
	cp -rf ./web/dist/static/* ./assets/tail/static/
	cp -rf ./web/dist/favicon.ico ./assets/tail/static/
	cp -rf ./web/dist/index.html ./assets/tail/static/
	go generate ./assets/...


log-tail:
	go build -o bin/log-tail ./cmd

clean:
	rm -f ./bin/log-tail

gotest:
	go test -v --cover -race -coverprofile=coverage.txt -covermode=atomic ./tests/ci/...


web-update: web-build file

web-build:
	cd web && npm run build
