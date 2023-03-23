FROM golang:1.20-alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

FROM golang:1.20-alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /neko
WORKDIR /neko
RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -ldflags="-s -w" -o /bin/neko ./cmd/app

FROM scratch
COPY --from=builder /neko/config /config
COPY --from=builder /neko/migrations /migrations
COPY --from=builder /bin/neko /app

ENTRYPOINT ["/app"]
