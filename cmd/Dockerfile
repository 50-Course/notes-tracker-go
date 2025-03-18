FROM golang:latest AS builder

WORKDIR /app

COPY . .
RUN go mod tidy

RUN go build -o grpc_service ./cmd/main.go

# RUN go build -o api_gateway ./api/gateway/server.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/grpc_service .
# COPY --from=builder /app/api_gateway .

EXPOSE 50051

CMD ["sh", "-c", "./grpc_service"]
