FROM golang:onbuild

ARG COMMIT
ENV COMMIT ${COMMIT}

EXPOSE 9000

ENTRYPOINT ["/go/bin/app", "-port", "9000", "-openbrowser", "false"]
