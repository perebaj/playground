FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
RUN CGO_ENABLED=0 go build -o compression-bench .

FROM alpine:3.20
COPY --from=builder /app/compression-bench /usr/local/bin/compression-bench
ENTRYPOINT ["compression-bench"]
