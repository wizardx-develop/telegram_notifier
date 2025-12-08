FROM golang:1.24.5 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w -s' -o /main

FROM alpine:latest
WORKDIR /
COPY --from=builder /main /main
ENTRYPOINT ["/main"]