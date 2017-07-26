FROM golang:alpine

RUN apk --no-cache add git

ADD main.go /go/src/app/main.go
ADD static /go/src/app/static
WORKDIR /go/src/app
RUN go get -d -v
RUN go install -v

ARG COMMIT
ENV COMMIT ${COMMIT}

EXPOSE 9000
CMD ["/go/bin/app", "-port", "9000"]
