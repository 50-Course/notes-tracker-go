# Dead Simple Task tracking application

### Architectural decisions:

- Follows a microservices architecture
- Uses a RESTful API for external communication and gRPC for internal communication
- Fully containerized using Docker
- And, Bootstrapped using Docker Compose

## System Architecture

```mermaid
graph TD;
    subgraph "Client"
        User
    end

    subgraph "API Gateway (BunRouter)"
        GatewayServer["API Gateway (BunRouter)"]
    end

    subgraph "Internal Service (gRPC Server)"
        gRPCServer["gRPC Server"]
    end

    subgraph "Database Layer"
        PostgreSQL["PostgreSQL (Bun ORM)"]
    end

    %% Arrows (Flow of Communication)
    User -->|REST API| GatewayServer
    GatewayServer -->|gRPC Call| gRPCServer
    gRPCServer -->|Database Queries| PostgreSQL
    PostgreSQL -->|Data Response| gRPCServer
    gRPCServer -->|gRPC Response| GatewayServer
    GatewayServer -->|HTTP Response| User
```

### **How It Works:**

- **User** makes HTTP requests to the **API Gateway**.
- **API Gateway** translates requests into **gRPC** calls to the **Internal Service**.
- **Internal Service** communicates with **PostgreSQL** using **Bun ORM**.
- Responses flow back **through the same pipeline** to the user.

## Installation

## Usage

## Contributing

## License
