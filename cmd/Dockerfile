FROM golang:1.21 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy && go build -o internal_service ./cmd/internal_service

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/internal_service .

CMD ["./internal_service"]
