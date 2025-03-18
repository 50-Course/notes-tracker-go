# Dead Simple Task tracking application

### Architectural decisions:

- Follows a microservices architecture
- Uses a RESTful API for external communication and gRPC for internal communication
- Fully containerized using Docker
- And, Bootstrapped using Docker Compose

## System Architecture

```mermaid
graph LR
    API_Gateway[API Gateway (HTTP Server)] -->|REST| BunRouter(BunRouter api/gateway/server.go);
    BunRouter -->|gRPC| Internal_Service[Internal Service (gRPC Server)];
    Internal_Service -->|ORM| PostgreSQL[PostgreSQL (Bun ORM)];
```

## Installation

## Usage

## Contributing

## License
