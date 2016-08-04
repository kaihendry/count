FROM golang
ADD . /go/src/github.com/spuul/count
RUN go get github.com/skratchdot/open-golang/open
RUN go install github.com/spuul/count
ARG COMMIT
ENV COMMIT ${COMMIT}
ENTRYPOINT ["/go/bin/count", "-port", "9000", "-openbrowser", "false"]
EXPOSE 9000
