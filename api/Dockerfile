# Image for our Todo HTTP API server
# Use multi-stage builds to reduce the size of the final image
FROM golang:latest AS builder

WORKDIR /app

# COPY go.mod go.sum ./
COPY . .
RUN go mod tidy


RUN go build -o api_gateway ./api/gateway/server.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/api_gateway .

EXPOSE 8080

CMD ["sh","-c", "./api_gateway"]

