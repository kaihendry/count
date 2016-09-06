FROM golang:1.7-alpine

ADD main.go /go/src/app/main.go
ADD static /go/src/app/static/

WORKDIR /go/src/app

RUN apk --no-cache add curl git && \
	go get -v -d && \
	go build && go install -v

ARG COMMIT
ENV COMMIT ${COMMIT}

EXPOSE 9000

ENTRYPOINT ["/go/bin/app", "-port", "9000", "-openbrowser", "false"]
