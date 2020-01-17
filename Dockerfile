FROM golang as builder

WORKDIR /app

COPY . .

ARG TARGET_ARCH=amd64
ARG VERSION

LABEL org.label-schema.version=$VERSION

RUN echo Building for ${TARGET_ARCH}
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGET_ARCH} \
	go build -ldflags "-X main.Version=${VERSION}"

FROM scratch
COPY --from=builder /app/count /app/
COPY --from=builder /app/static /app/static
ENV PORT 8080

WORKDIR /app
ENTRYPOINT ["/app/count"]
