all: clean fmt build

build: log-tail

fmt:
	go fmt ./...

log-tail:
	go build -o bin/log-tail ./cmd

clean:
	rm -f ./bin/log-tail


web-build:
	cd web && npm run build

web-deploy:
	rm -rf assets/tail/static && \
		mv web/dist/static assets/tail && \
		mv web/dist/favicon.ico assets/tail/static && mv web/dist/index.html assets/tail/static && \
		rm -rf web/dist