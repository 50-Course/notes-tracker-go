# Image for our Todo HTTP API server
# Use multi-stage builds to reduce the size of the final image
FROM golang:latest AS builder

WORKDIR /app

# COPY . .
COPY go.mod go.sum ./

RUN go mod tidy
COPY ./config/.env /app/config/.env

COPY . .


RUN go build -o api_gateway ./api/gateway/server.go

FROM ubuntu:22.04
WORKDIR /root/
COPY --from=builder /app/api_gateway .
COPY --from=builder /app/config/.env /root/config/.env
RUN chmod +x /root/api_gateway

EXPOSE 8080

CMD ["./api_gateway"]

