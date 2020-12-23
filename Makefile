VERSION = $(shell git rev-parse --abbrev-ref HEAD)-$(shell git rev-parse --short HEAD)

.PHONY: build clean deploy

count:
	env CGO_ENABLED=0 GOOS=linux go build  -ldflags="-s -w -X main.Version=${VERSION}" -o count *.go

clean:
	rm -f count

deploy: count
	sls deploy --verbose
