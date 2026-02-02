# Go Backend Boilerplate

A production-ready Go backend boilerplate following clean architecture principles, featuring dependency injection, modular monolith structure, and comprehensive observability.

## 🚀 Features

- **Clean Architecture**: Modular monolith structure with clear separation of concerns (Handler -> Usecase -> Repository).
- **Dependency Injection**: Uses `github.com/samber/do` for a type-safe, lightweight DI container.
- **Authentication**: Full JWT implementation with Access and Refresh tokens. Supports token rotation and revocation.
- **Database**: GORM with PostgreSQL, featuring configurable connection pooling and auto-migrations.
- **Observability**: 
  - Structured Logging with `log/slog`.
  - Standard Health Checks (`/health/live`, `/health/ready`).
- **Resilience**: 
  - Rate Limiting.
  - Graceful Shutdown handling OS signals (SIGINT, SIGTERM).
  - Recovery middleware.
- **Validation**: Request body validation using `go-playground/validator`.
- **Testing**: Unit test infrastructure with `testify` and mocks.

## 🛠 Tech Stack

- **Go** (1.21+)
- **Echo v5** (Web Framework)
- **GORM** (ORM)
- **PostgreSQL** (Database)
- **Slog** (Structured Logging)
- **Cleanenv** (Configuration Management)
- **JWT-Go** (Authentication)

## 🏗 Project Structure

```text
.
├── cmd/
│   └── api/             # Main entry point
├── internal/
│   ├── auth/            # Auth domain (Login, Refresh, Logout)
│   │   ├── entity.go
│   │   ├── repository.go
│   │   ├── usecase.go
│   │   ├── handler.go
│   │   └── setup.go
│   ├── config/          # Configuration management
│   ├── crypto/         # Crypto domain (example)
│   │   ├── portfolio/   # Portfolio subdomain
│   │   │   ├── entity.go
│   │   │   ├── repository.go
│   │   │   ├── usecase.go
│   │   │   └── handler.go
│   │   └── setup.go     # Domain DI setup
│   ├── database/        # Database connection and helpers
│   ├── infra/
│   │   ├── auth/        # Infrastructure level auth (JWT Service, Middleware)
│   │   └── health/      # Health check probes
│   └── router/          # Echo router and middleware configuration
└── pkg/
    └── response/        # Standardized API response helpers
```

## 🚦 Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/excflor/go-boilerplate.git
   cd go-boilerplate
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Setup environment variables:
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials and JWT secret
   ```

### Running the App

- **Development (with live reload)**:
  ```bash
  air
  ```
- **Standard Run**:
  ```bash
  go run ./cmd/api/main.go
  ```

## 🔐 Authentication

The project uses JWT for authentication with support for refresh tokens.

> [!NOTE]
> For demonstration purposes in this boilerplate, the **login credentials are mocked**. 
> - **Username**: `admin`
> - **Password**: `admin`
> - **Mock User ID**: `7ea078fa-aac0-4364-8f5f-ba69b136b8f7`

1. **Login**: Get tokens by calling `POST /auth/login`.
   ```bash
   curl -X POST http://localhost:4001/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username": "admin", "password": "admin"}'
   ```
2. **Refresh**: Get a new access token using your refresh token.
   ```bash
   curl -X POST http://localhost:4001/auth/refresh \
     -H "Content-Type: application/json" \
     -d '{"refresh_token": "<your_refresh_token>"}'
   ```
3. **Logout**: Revoke your refresh token.
   ```bash
   curl -X POST http://localhost:4001/auth/logout \
     -H "Content-Type: application/json" \
     -d '{"refresh_token": "<your_refresh_token>"}'
   ```
4. **Authorize**: Add the access token to the `Authorization` header for protected routes:
   ```text
   Authorization: Bearer <your_access_token>
   ```

## 🏥 Health Checks

- **Liveness**: `GET /health/live` (Is the process running?)
- **Readiness**: `GET /health/ready` (Is the DB connected and ready?)

## 🧪 Testing

Run all tests:
```bash
go test ./...
```