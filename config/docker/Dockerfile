FROM golang:1.21-alpine AS builder

WORKDIR /build
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o kube-ipam ./cmd

FROM scratch
COPY --from=builder /build/kube-ipam /kube-ipam
ENTRYPOINT ["/kube-ipam"]