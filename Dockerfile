FROM golang:1.27 AS builder

WORKDIR /src
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY pkg ./pkg
COPY cmd ./cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o manager ./cmd/manager

FROM ubuntu:latest
WORKDIR /
COPY --from=builder /src/manager /manager
ENTRYPOINT ["/manager"]
