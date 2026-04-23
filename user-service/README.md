# User Service

## Overview
User Service is a microservice responsible for user management and authentication within the Booking System platform. It provides robust features for user registration, login, profile management, and secure access control using JWT and API Key validation.

## Tech Stack
- **Language:** Go 1.25.0
- **Web Framework:** [Gin Gonic](https://gin-gonic.com/)
- **ORM:** [GORM](https://gorm.io/)
- **Database:** PostgreSQL
- **Authentication:** JWT (JSON Web Token) & Bcrypt for password hashing
- **Configuration:** [Viper](https://github.com/spf13/viper) & [Godotenv](https://github.com/joho/godotenv)
- **CLI:** [Cobra](https://github.com/spf13/cobra)
- **Rate Limiting:** [Tollbooth](https://github.com/didip/tollbooth)
- **Logging:** [Logrus](https://github.com/sirupsen/logrus)
- **Validation:** [Go Playground Validator](https://github.com/go-playground/validator)

## Project Structure
The project follows a layered architecture to ensure separation of concerns and maintainability:

```text
├── cmd/                # Application entry point and CLI commands
├── common/             # Shared utilities (standardized responses and error handling)
├── config/             # Configuration loading (Viper) and DB initialization
├── constants/          # Global constants and error definitions
├── controllers/        # Request handlers (input validation and response formatting)
├── database/           # Database migrations and seeders
├── domain/             # Data models (GORM) and DTOs (Data Transfer Objects)
├── middleware/         # Custom Gin middlewares (Auth, Rate Limiting, Panic Recovery)
├── repositories/       # Data access layer (DB queries)
├── routes/             # Route definitions and registration
├── services/           # Business logic layer
└── main.go             # Main entry point that executes Cobra commands
```

## Prerequisites
- **Go:** 1.25 or higher
- **PostgreSQL:** A running instance for data storage
- **Air:** (Optional) For live-reloading during development

## Installation & Setup
1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd user-service
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Environment Configuration:**
   Copy the example environment file and update the values:
   ```bash
   cp .env.example .env
   ```

4. **Database Setup:**
   Ensure your PostgreSQL database is running and the credentials in `.env` are correct. The application will automatically run migrations on startup.

## How to Run
### Development Mode (with Air)
If you have [Air](https://github.com/air-verse/air) installed:
```bash
air
```

### Manual Run
```bash
go run main.go serve
```

### Production Build
```bash
go build -o user-service .
./user-service serve
```

## API Documentation

### Base URL: `/api/v1`

#### Public Endpoints
| Method | Path | Description |
| :--- | :--- | :--- |
| `GET` | `/` | Health check endpoint |
| `POST` | `/auth/register` | Register a new user |
| `POST` | `/auth/login` | Authenticate user and receive JWT |

#### Protected Endpoints
*Required Headers: `Authorization: Bearer <token>`, `x-api-key`, `x-request-at`, `x-service-name`*

| Method | Path | Description |
| :--- | :--- | :--- |
| `GET` | `/auth/user` | Get profile of the currently logged-in user |
| `GET` | `/auth/user/:uuid` | Get profile of a specific user by UUID |
| `PUT` | `/auth/:uuid` | Update user profile |

### API Key Validation (SHA256)
For protected endpoints, an API key is required. It is generated using:
`SHA256(serviceName + ":" + signatureKey + ":" + requestAt)`

In **development mode**, you can use the helper endpoint to generate a valid key:
`GET /dev/api-key?service_name=my-service&request_at=2026-04-18T12:00:00Z`

## Error Codes
The service uses standardized HTTP status codes and custom error messages:

| Error Message | HTTP Status | Meaning |
| :--- | :--- | :--- |
| `internal server error` | 500 | Unexpected server error |
| `unauthorized` | 401 | Missing or invalid token/API key |
| `forbidden` | 403 | Insufficient permissions |
| `too many requests` | 429 | Rate limit exceeded |
| `user not found` | 404 | User does not exist |
| `username already exists` | 409 | Conflict during registration |
| `email already exists` | 409 | Conflict during registration |
| `username or password is incorrect`| 401 | Authentication failure |
