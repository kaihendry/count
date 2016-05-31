FROM golang
RUN go get -x github.com/kaihendry/count
WORKDIR /go/src/github.com/kaihendry/count/
ENTRYPOINT /go/bin/count -port 9000 -openbrowser false
EXPOSE 9000
