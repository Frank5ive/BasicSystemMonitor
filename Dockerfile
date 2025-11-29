# Build stage
FROM golang:1.24-bullseye AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags='-s -w' -o /out/monitor ./

# Final image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /out/monitor /usr/local/bin/monitor
USER 1000:1000
ENTRYPOINT ["/usr/local/bin/monitor"]
