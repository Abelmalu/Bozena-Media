# Golang Posts API

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)
![Gin](https://img.shields.io/badge/Gin-Framework-00ADD8?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-v15+-336791?style=for-the-badge&logo=postgresql)

A robust backend REST API for a social media posting platform, built with Go and the Gin framework. This project demonstrates clean architecture patterns, secure authentication with JWT, and role-based access control.

## 🚀 Features

- **Authentication & Authorization**
  - specific User Registration and Login handles
  - JWT-based stateless authentication
  - Token Refresh mechanism
  - Secure Logout
  - Role-Based Access Control (RBAC) middleware
- **Post Management**
  - Create and viewing posts
  - Owner-only Edit and Delete (granular permissions)
- **Technical Highlights**
  - **Framework**: [Gin Web Framework](https://gin-gonic.com/) for high performance.
  - **Database**: PostgreSQL with `pgx` driver.
  - **Live Reload**: Configured with [Air](https://github.com/air-verse/air) for rapid development.
  - **Containerization**: Docker support included.

## 🛠️ Project Structure

```bash
├── cmd/
│   └── main.go           # Application entry point
├── internal/
│   ├── auth/             # Authentication logic (Register, Login, Refresh)
│   ├── middleware/       # Auth & Role-based middleware
│   ├── posts/            # Post CRUD operations
│   ├── models/           # Data models
│   └── routes.go         # Router configuration
├── pkg/                  # Shared packages (DB connection)
├── migrations/           # Database migrations
└── Dockerfile            # Docker configuration
```

## ⚡ Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.25+)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Air](https://github.com/air-verse/air) (Optional, for hot reload)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/abelmalu/golang-posts.git
   cd golang-posts
   ```

2. **Environment Setup**
   Create a `.env` file in the root directory (refer to your config or code for required variables, typically `DB_URL`, `JWT_SECRET`, etc.).

   ```bash
   touch .env
   ```

3. **Install Dependencies**
   ```bash
   go mod download
   ```

4. **Run the Application**

   **Using Air (Recommended for Dev):**
   ```bash
   air
   ```

   **Using Go:**
   ```bash
   go run cmd/main.go
   ```

   The server will start on `http://localhost:8080`.

## 🔌 API Endpoints

### Authentication

| Method | Endpoint    | Description                          | Auth Required |
| :----- | :---------- | :----------------------------------- | :------------ |
| `POST` | `/register` | Register a new user                  | ❌            |
| `POST` | `/login`    | Login and receive Access/Refresh JWT | ❌            |
| `POST` | `/refresh`  | Refresh expired access token         | ❌            |
| `POST` | `/logout`   | Invalidate session                   | ✅            |

### Posts

*Note: All Post endpoints require the `users` role.*

| Method   | Endpoint     | Description            | Permissions          |
| :------- | :----------- | :--------------------- | :------------------- |
| `GET`    | `/posts`     | List all posts         | Authenticated User   |
| `POST`   | `/posts`     | Create a new post      | Authenticated User   |
| `PUT`    | `/posts/:id` | Update a specific post | **Owner Only**       |
| `DELETE` | `/posts/:id` | Delete a specific post | **Owner Only**       |

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
