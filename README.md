# Dead Simple Task tracking application

### Architectural decisions:

- Follows a microservices architecture
- Uses a RESTful API for external communication and gRPC for internal communication
- Fully containerized using Docker
- And, Bootstrapped using Docker Compose

### System Architecture

```mermaid
    A[API Gateway (HTTP Server)] -->|REST| B(BunRouter api/gateway/server.go);
    B -->|gRPC| C[Internal Service (gRPC Server)];
    C -->|ORM| D[PostgreSQL (Bun ORM)];
```

## Installation

## Usage

## Contributing

## License
