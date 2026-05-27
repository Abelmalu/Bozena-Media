# 🌐 Golang Social Media Microservices

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)
![gRPC](https://img.shields.io/badge/gRPC-Framework-4285F4?style=for-the-badge&logo=grpc)
![Gin](https://img.shields.io/badge/Gin-Framework-00ADD8?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-v15+-336791?style=for-the-badge&logo=postgresql)
![Redis](https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white)

A modern, scalable social network backend built using **Microservices Architecture**. This project leverages **gRPC** for efficient internal communication between services and **Gin** as the entry point via an API Gateway.

---

## 🏗️ Architecture Overview

The project is built using **Clean Architecture** principles in each service to ensure scalability, testability, and a clear separation of concerns.

### 🏛️ Clean Architecture Layers
Each service (`Auth`, `post`,`like`) is structured into the following layers within the `internal/` directory:

- **`internal/handlers/`**: The **Delivery/Transport Layer**. It implements the gRPC server interfaces and handles incoming requests and outgoing responses.
- **`internal/service/`**: The **Business Logic/Use Case Layer**. This contains the core logic of the application and is independent of external frameworks.
- **`internal/repository/`**: The **Data Access/Infrastructure Layer**. It handles all database operations (PostgreSQL) and isolates the data storage details from the business logic.
- **`internal/core/`**: The **Entities/Domain Layer**. Contains core business rules and logic that are essential to the service.
- **`internal/models/`**: The **Data Models**. Defines the structures used for data transfer and storage across different layers.

### 🛰️ The Services

1.  **API Gateway**: The entry point for all client requests. It handles routing and authenticates requests using the **Auth Service** before forwarding them to the internal services.
2.  **Auth Service**: Manages user registration, login, JWT issuance, and session management.
3.  **Post Service**: Handles all post-related operations (CRUD).
4.  **Likes Service**: Manages post reactions and likes.
5.  **Feeds Service** *(Coming Soon 🚀)*: Will generate user-specific timelines and feeds.

### 🔌 Communication Map

- **External (Client -> Gateway)**: REST API (HTTP/JSON)
- **Internal (Gateway -> Services)**: gRPC (Protobuf)
- **Internal (Service -> Service)**: gRPC (Protobuf)

---

## 🚀 Features

- **Clean Architecture Implementation**: Strict separation of concerns for maintainability and testability.
- **Microservices Architecture**: Decoupled services for better scalability and maintenance.
- **gRPC Integration**: High-performance internal communication using Protocol Buffers.
- **Secure Authentication**: JWT-based stateless authentication with token refresh, secure logout using **Redis** for token blacklisting, and cross-platform support (HttpOnly cookies for Web, JSON response for Mobile).
- **Distributed Rate Limiting**: Token bucket rate limiting applied to routes, leveraging **Redis** for distributed state management and **Lua scripts** to guarantee atomicity and prevent race conditions.
- **Structured Error Handling**: Custom error mapper (`ierrors`) that safely translates internal gRPC errors into clean, user-friendly HTTP responses while preserving internal causes for debugging.
- **Request Tracing**: Distributed tracing using `RequestID` propagation across the API Gateway and microservices for easier log correlation.
- **Structured Logging**: High-performance, structured JSON logging implemented using **Uber's Zap** (`go.uber.org/zap`).
- **Data Validation**: Strict incoming request validation utilizing **go-playground/validator** and Gin's built-in bindings.
- **Granular Permissions**: Role-Based Access Control (RBAC) and ownership verification.
- **Database Migrations**: Managed PostgreSQL schemas for each service.
- **Containerization**: Docker support for simplified deployment.

---

## 🛠️ Project Structure

```bash
├── APIGateway/           # Entry point (Gin + gRPC Clients)
├── Auth/                 # Authentication service (gRPC Server)
├── post/                 # Post management service (gRPC Server)
├── like/                 # Like management service (gRPC Server)
├── pkg/                  # Shared utilities (JWT, etc.)
├── migrations/           # Database migration files
└── proto/                # Shared Protocol Buffer definitions (if applicable)
```

---

## ⚡ Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.25+)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Redis](https://redis.io/download/) (for distributed rate limiting and token blacklisting)
- [Protoc](https://grpc.io/docs/protoc-installation/) (for generating gRPC code)

### Installation & Running

1. **Clone the repository**
   ```bash
   git clone https://github.com/abelmalu/golang-posts.git
   cd golang-posts
   ```

2. **Setup Environment Variables**
   Each service (`APIGateway`, `Auth`, `post`) requires its own `.env` file. Refer to the `.env.example` in each directory.

3. **Run the Services** (Open separate terminals):

   **Auth Service (Port 50052):**
   ```bash
   cd Auth
   go run cmd/main.go
   ```

   **Post Service (Port 50051):**
   ```bash
   cd post
   go run cmd/main.go
   ```

   **Like Service (Port 50053):**
   ```bash
   cd like
   go run cmd/main.go
   ```

   **API Gateway (Port 8080):**
   ```bash
   cd APIGateway
   go run cmd/gateway/main.go
   ```

---

## 🔌 API Endpoints (via Gateway)

### Authentication

| Method | Endpoint              | Description                          | Auth Required |
| :----- | :-------------------- | :----------------------------------- | :------------ |
| `POST` | `/api/auth/register`  | Register a new user                  | ❌            |
| `POST` | `/api/auth/login`     | Login and receive Access/Refresh JWT | ❌            |
| `POST` | `/api/auth/refresh`   | Refresh expired access token         | ❌            |
| `POST` | `/api/auth/logout`    | Invalidate session                   | ✅            |

### Posts

| Method   | Endpoint           | Description            | Permissions          |
| :------- | :----------------- | :--------------------- | :------------------- |
| `GET`    | `/api/posts/`      | List all posts         | Authenticated User   |
| `POST`   | `/api/posts/`      | Create a new post      | Authenticated User   |
| `PUT`    | `/api/posts/:id`   | Update a specific post | **Owner Only**       |
| `DELETE` | `/api/posts/:id`   | Delete a specific post | **Owner Only**       |

### Likes

| Method   | Endpoint              | Description                  | Permissions          |
| :------- | :-------------------- | :--------------------------- | :------------------- |
| `POST`   | `/api/posts/like/:id` | Toggle like/unlike on a post | Authenticated User   |

---

## 📅 Roadmap

- [x] **Likes Service**: Allow users to like and react to posts.
- [ ] **Feeds Service**: Optimized algorithm for displaying posts.
- [ ] **Docker Compose**: Single command to spin up the entire ecosystem.
- [ ] **Service Discovery**: Implement Consul or Etcd.

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
