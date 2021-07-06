# Build the manager binary
FROM golang:1.15 as builder

# Copy in the go src
RUN ls -la
WORKDIR /go/src/github.com/infinimesh/enterprise-operator
COPY pkg/    pkg/
COPY cmd/    cmd/
COPY vendor/ vendor/
COPY pkg/ vendor/github.com/infinimesh/enterprise-operator/pkg/
RUN ls -laR pkg/ vendor/github.com/googleapis/
RUN ls -laR go/ /go/
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager github.com/infinimesh/enterprise-operator/cmd/manager

# Copy the controller-manager into a thin image
FROM ubuntu:latest
WORKDIR /
COPY --from=builder /go/src/github.com/infinimesh/enterprise-operator/manager .
ENTRYPOINT ["/manager"]
