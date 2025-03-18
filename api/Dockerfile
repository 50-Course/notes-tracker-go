# Image for our Todo HTTP API server
# Use multi-stage builds to reduce the size of the final image
FROM golang:1.24-alphine AS builder

WORKDIR /app

COPY . .
RUN go mod tidy && go build -o gateway ./cmd/gateway/

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/gateway .

CMD ["./gateway"]

