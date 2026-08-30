# 🌐 Golang Social Media Microservices

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)
![gRPC](https://img.shields.io/badge/gRPC-Framework-4285F4?style=for-the-badge&logo=grpc)
![Gin](https://img.shields.io/badge/Gin-Framework-00ADD8?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-v15+-336791?style=for-the-badge&logo=postgresql)
![MongoDB](https://img.shields.io/badge/MongoDB-%234ea94b.svg?style=for-the-badge&logo=mongodb&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white)
![Kafka](https://img.shields.io/badge/Apache_Kafka-231F20?style=for-the-badge&logo=apachekafka&logoColor=white)

**Bozena Media** is a distributed, high-performance social media backend platform built with Go using a microservices architecture. Architected to support high-throughput real-time interactions, it integrates **gRPC** for low-latency synchronous communication, **Apache Kafka** for event-driven asynchronous streaming, and **WebSockets/SSE** for instantaneous message delivery and live notifications.

---

## 🏗️ Architecture Overview

The system is architected around **Clean Architecture** principles to establish a robust separation of concerns, testability, and scalability across all backend microservices. Each microservice (`Auth`, `post`, `like`, `follow`, `Feed`, `notification`, `Chat`) is decoupled, maintaining its own database and models. 

Internal communication uses a hybrid topology:
- **Synchronous RPC**: gRPC (via Protocol Buffers) provides high-performance, low-latency inter-service queries (e.g. Gateway to Auth/Post).
- **Asynchronous Events**: Apache Kafka facilitates event-driven data propagation (e.g. `userCreated`, `postCreated`, `liked`, `followed` events) to keep decentralized read caches eventually consistent across the services.

### 🏛️ Clean Architecture Layers
Each service (`Auth`, `post`, `like`, `follow`, `feed`, `notification`, `Chat`) is structured into the following layers within the `internal/` directory:

- **`internal/handlers/`**: The **Delivery/Transport Layer**. It implements the gRPC/HTTP server interfaces and WebSocket handlers, managing incoming requests and outgoing responses.
- **`internal/service/`**: The **Business Logic/Use Case Layer**. This contains the core logic of the application and is independent of external frameworks.
- **`internal/repository/`**: The **Data Access/Infrastructure Layer**. It handles database operations (PostgreSQL / MongoDB) and isolates the data storage details from the business logic.
- **`internal/core/`**: The **Entities/Domain Layer**. Contains core business rules and logic that are essential to the service.
- **`internal/models/`**: The **Data Models**. Defines the structures used for data transfer and storage across different layers.

### 🛰️ The Services

1.  **API Gateway**: The public entry point. It handles routing, request IDs, CORS, rate limiting, auth checks, and forwards/proxies requests to internal services over gRPC or HTTP/WebSocket.
2.  **Auth Service**: Manages user registration, login, refresh, logout, JWT issuance, session management, and username search.
3.  **Post Service**: Handles post CRUD and user-specific post lookup.
4.  **Like Service**: Manages like and unlike operations plus post like listings.
5.  **Follow Service**: Manages follow and unfollow relationships between users, plus follower/following listings.
6.  **Feed Service**: Builds the authenticated user feed from post and user metadata.
7.  **Notification Service**: Connects to Kafka and pushes real-time follower notifications to the browser via Server-Sent Events (SSE).
8.  **Chat Service**: Manages direct messaging between users. It uses **MongoDB** to persist conversations, direct messages, and user cache info, and supports real-time communication via WebSockets.

### 🔌 Communication Map

- **External (Client -> Gateway)**: REST API (HTTP/JSON) & Server-Sent Events (SSE)
- **Internal (Gateway -> Services)**: gRPC (Protobuf)
- **Internal (Service -> Service)**: gRPC and kafka

### 🔁 Request Flow

1. A client calls the API Gateway.
2. The Gateway authenticates the request, applies middleware, and attaches request metadata.
3. The Gateway forwards the request to the relevant service over gRPC.
4. Services persist data in PostgreSQL, use Redis for auth/rate-limit state, and emit Kafka events where needed.

### 🗂️ Static File Storage

- Static files such as user avatars and post images are stored in MinIO object storage.
- Services generate presigned URLs so clients can upload and access media securely without exposing the bucket directly.

### 🧩 Data Ownership

- Each service owns its own database (PostgreSQL for Auth, Post, Like, Follow, Feed, Notification; MongoDB for Chat) and migrations/collections.
- Auth, Post, Follow, Feed, and Chat all maintain user caches for username/name/avatar lookups.
- Kafka propagates `userCreated` and post-related events to downstream services, and triggers real-time pushes in the Notification Service.
- Redis is used for refresh-token/session blacklist state and distributed token-bucket rate limiting.

---

## 🚀 Features

- **Clean Architecture Implementation**: Strict separation of concerns for maintainability and testability.
- **Microservices Architecture**: Decoupled services for better scalability and maintenance.
- **Real-Time Direct Messaging**: Handled by the **Chat Service** using **Gorilla WebSockets** for low-latency, real-time message streaming.
- **NoSQL Chat Persistence**: Conversation history and direct messages are saved to **MongoDB** with cursor-based pagination for high write throughput and schema flexibility.
- **Real-Time Notifications**: Server-Sent Events (SSE) stream real-time updates (like new followers) from the Notification Service directly to the frontend, proxied securely by the API Gateway.
- **Event-Driven Architecture**: Configured with **Apache Kafka**. When users register or update their profiles, `userCreated` events are published and consumed by other services. This asynchronous messaging is also utilized for post creation logic.
- **gRPC Integration**: High-performance internal communication using Protocol Buffers.
- **Secure Authentication**: JWT-based stateless authentication with token refresh, secure logout using **Redis** for token blacklisting, and cross-platform support (HttpOnly cookies for Web, JSON response for Mobile).
- **Distributed Rate Limiting**: Token bucket rate limiting applied to routes, leveraging **Redis** for distributed state management and **Lua scripts** to guarantee atomicity and prevent race conditions.
- **Structured Error Handling**: Custom error mapper (`ierrors`) that safely translates internal gRPC errors into clean, user-friendly HTTP responses while preserving internal causes for debugging.
- **Request Tracing**: Distributed tracing using `RequestID` propagation across the API Gateway and microservices for easier log correlation.
- **Structured Logging**: High-performance, structured JSON logging implemented using **Uber's Zap** (`go.uber.org/zap`).
- **Data Validation**: Strict incoming request validation utilizing **go-playground/validator** and Gin's built-in bindings.
- **Granular Permissions**: Role-Based Access Control (RBAC) and ownership verification.
- **Independent Database & Migrations**: Each service fully owns its own database and database migrations, ensuring strong data decoupling across the microservices ecosystem.
- **Containerization**: Docker support for simplified deployment.

---

## 🛠️ Project Structure

```bash
├── APIGateway/           # Entry point (Gin + gRPC Clients)
├── Auth/                 # Authentication service (gRPC Server)
├── post/                 # Post management service (gRPC Server)
├── like/                 # Like management service (gRPC Server)
├── follow/               # Follow management service (gRPC Server)
├── feed/                 # Feed management service (gRPC Server)
├── notification/         # Real-time notifications service (HTTP SSE Server)
├── Chat/                 # Real-time Chat service (HTTP/WebSocket server using MongoDB)
├── frontend/             # React + TypeScript web client
├── pkg/                  # Shared utilities (JWT, etc.)
├── migrations/           # Database migration files
└── proto/                # Shared Protocol Buffer definitions (if applicable)
```

---

## ⚡ Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.25+)
- [PostgreSQL](https://www.postgresql.org/download/)
- [MongoDB](https://www.mongodb.com/try/download/community) (for Chat history and user cache)
- [Redis](https://redis.io/download/) (for distributed rate limiting and token blacklisting)
- [Apache Kafka](https://kafka.apache.org/downloads) (for event-driven communication)
- [MinIO](https://min.io/) (for object storage of user avatars and post media)
- [Protoc](https://grpc.io/docs/protoc-installation/) (for generating gRPC code)
- [Air](https://github.com/air-verse/air) (optional, for live-reload development)

### Installation & Running

1. **Clone the repository**
   ```bash
   git clone https://github.com/abelmalu/golang-posts.git
   cd golang-posts
   ```

2. **Setup Environment Variables (For Local Execution)**
   If running services locally directly on your host machine, each service requires its own `.env` file containing local credentials. Refer to the `.env.example` file located in each service directory.

3. **Run the Services** (Open separate terminals):
   You can run each service using the standard Go command or use **Air** for live-reloading during development (recommended, as `.air.toml` config files are provided in each directory).

   **Auth Service (Port 50052):**
   ```bash
   cd Auth
   # Using Air (hot reload)
   air
   # Or using standard Go run
   go run cmd/main.go
   ```

   **Post Service (Port 50051):**
   ```bash
   cd post
   air # or: go run cmd/main.go
   ```

   **Like Service (Port 50053):**
   ```bash
   cd like
   air # or: go run cmd/main.go
   ```

   **Follow Service (Port 50054):**
   ```bash
   cd follow
   air # or: go run cmd/main.go
   ```

   **Feed Service (Port 50055):**
   ```bash
   cd Feed
   air # or: go run cmd/main.go
   ```

   **Notification Service (Port 8083 / Stream SSE):**
   ```bash
   cd notification
   air # or: go run cmd/main.go
   ```

   **Chat Service (Port 8084 / WebSockets):**
   ```bash
   cd Chat
   air # or: go run cmd/main.go
   ```

   **API Gateway (Port 8080):**
   ```bash
   cd APIGateway
   air -c .air.toml # or: go run cmd/gateway/main.go
   ```

### Running with Docker Compose 🐳

Docker Compose provides a fully self-contained environment to run the entire backend ecosystem with a single command. **There is no need to manually configure `.env` files** when running with Docker Compose, as all necessary configurations and environment variables are already baked into `compose.yml`.

To spin up the entire environment:

1. Make sure Docker and the Docker Compose plugin are installed on your system.
2. Run the stack from the root directory:
   ```bash
   docker compose up --build
   ```

When you execute this command, Docker Compose will launch and configure:
* **Databases & Cache**: PostgreSQL (relational DB), MongoDB (chat storage), Redis (session blacklist & rate limiting), MinIO (object storage).
* **Message Broker**: Apache Kafka (including node controller and bootstrap setup).
* **Migration Runners**: Dedicated migration runner containers (`auth-migration`, `post-migration`, `like-migration`, `follow-migration`, `feed-migration`, `notification-migration` using `migrate/migrate:v4.17.0`) will automatically apply all database schemas/migrations before their respective services launch.
* **Microservices**: All backend service containers (`auth-service`, `post-service`, `like-service`, `follow-service`, `feed-service`, `notification-service`, `chat-service`, and `api-gateway`) automatically starting and wiring together.

Once running, **all microservice endpoints are exposed through the API Gateway** mapped to port `8080`. You can access and interact with the entire backend system locally via `http://localhost:8080` (e.g. testing routes, consuming Server-Sent Events, or establishing WebSocket connections).

---

## 🧪 Testing

The platform features unit tests covering the core business logic layer of the major microservices to ensure reliability and maintainability.

* **Table-Driven Tests**: Service layer tests are implemented using Go's table-driven pattern to test multiple success and error scenarios comprehensively.
* **Mocking**: External components (like repository, cache, and broker layers) are mocked using test interfaces to isolate the service layer logic.
* **Running Tests**: Run tests for all services from the workspace root:
  ```bash
  go test ./...
  ```
  Or test a specific service:
  ```bash
  cd Auth && go test ./...
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
| `GET`  | `/api/auth/search`    | Search users by username             | ✅            |

### Posts

| Method   | Endpoint                | Description            | Permissions          |
| :------- | :---------------------- | :--------------------- | :------------------- |
| `GET`    | `/api/posts/`           | List all posts         | Authenticated User   |
| `POST`   | `/api/posts/`           | Create a new post      | Authenticated User   |
| `GET`    | `/api/posts/user/:id`   | Get posts by a user    | Authenticated User   |
| `PUT`    | `/api/posts/update/:id` | Update a specific post | **Owner Only**       |
| `DELETE` | `/api/posts/delete/:id` | Delete a specific post | **Owner Only**       |

### Likes

| Method   | Endpoint              | Description                  | Permissions          |
| :------- | :-------------------- | :--------------------------- | :------------------- |
| `POST`   | `/api/posts/like/:id` | Toggle like/unlike on a post | Authenticated User   |
| `GET`    | `/api/posts/likes/:id`| Get total likes for a post   | Authenticated User   |

### Follows

| Method   | Endpoint                     | Description                  | Permissions          |
| :------- | :--------------------------- | :--------------------------- | :------------------- |
| `POST`   | `/api/follow/:id`            | Toggle follow/unfollow user  | Authenticated User   |
| `GET`    | `/api/follow/followers/:id`  | View followers of a user     | Authenticated User   |
| `GET`    | `/api/follow/followings/:id` | View users a user follows    | Authenticated User   |

### Feeds

| Method   | Endpoint              | Description                  | Permissions          |
| :------- | :-------------------- | :--------------------------- | :------------------- |
| `GET`    | `/api/feed/`          | Get user timeline feed       | Authenticated User   |

### Notifications

| Method | Endpoint                    | Description                           | Permissions          |
| :----- | :-------------------------- | :------------------------------------ | :------------------- |
| `GET`  | `/api/notifications/stream` | Real-time SSE stream for new followers| Authenticated User   |
| `GET`  | `/api/notification/user`    | Get past notifications for a user     | Authenticated User   |

### Chat

| Method | Endpoint                     | Description                           | Permissions          |
| :----- | :--------------------------- | :------------------------------------ | :------------------- |
| `GET`  | `/api/chat/ws`               | Establish real-time WebSocket session | Authenticated User   |
| `GET`  | `/api/chat/user/chats`       | Get direct message history & list     | Authenticated User   |
| `GET`  | `/api/chat/:id/messages`     | Get messages in a specific chat       | Authenticated User   |

---

## 📅 Roadmap

- [x] **Likes Service**: Allow users to like and react to posts.
- [x] **Follows Service**: Allow users to follow and unfollow other users.
- [x] **Feeds Service**: Optimized algorithm for displaying posts.
- [x] **Notification Service**: Real-time SSE pushes proxied via API Gateway.
- [x] **Chat Service**: Real-time direct messaging backed by **MongoDB**.
- [x] **Docker Compose**: Single command to spin up the entire ecosystem.
- [ ] **Service Discovery**: Implement Consul or Etcd.

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Frontend

The React frontend  which is vibe coded lives in [`frontend/`](/home/abel/Projects/GO/Bozena-Media/frontend) and talks only to the API Gateway.
