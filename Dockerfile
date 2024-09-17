FROM golang:1.17 as builder
WORKDIR /app
COPY . .
RUN go build -o distributed-file-system cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/distributed-file-system .
CMD ["./distributed-file-system"]
