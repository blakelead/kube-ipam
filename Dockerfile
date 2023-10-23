FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
COPY cmd cmd
COPY internal internal
COPY license .
RUN go mod tidy && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o kube-ipam /build/cmd

FROM scratch
COPY --from=builder /build/kube-ipam /kube-ipam
ENTRYPOINT ["/kube-ipam"]