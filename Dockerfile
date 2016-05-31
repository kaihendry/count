FROM golang
RUN go install github.com/kaihendry/count
ENTRYPOINT /go/bin/count -port 9000 -openbrowser false
EXPOSE 9000
