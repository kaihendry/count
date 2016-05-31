FROM golang
RUN go get -x github.com/kaihendry/count
ADD /static /go/src/github.com/kaihendry/count/static
ENTRYPOINT /go/bin/count -port 9000 -openbrowser false
EXPOSE 9000
