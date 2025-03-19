FROM golang:latest AS builder

WORKDIR /app

# COPY . .
COPY go.mod go.sum ./

RUN go mod tidy

COPY ./config/.env /app/config/.env

COPY . .

RUN go build -o grpc_service ./cmd/main.go

# RUN go build -o api_gateway ./api/gateway/server.go

FROM ubuntu:22.04
WORKDIR /root/
COPY --from=builder /app/grpc_service .
COPY --from=builder /app/config/.env /root/config/.env
RUN chmod +x /root/grpc_service

EXPOSE 50051

CMD ["./grpc_service"]
