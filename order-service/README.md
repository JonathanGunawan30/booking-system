# Order Service

A high-performance microservice built with Go for managing orders in a booking system. This service handles order creation, status tracking, and integrates with Payment and Field services via REST and Kafka.

## Features

- **Order Management**: Create and retrieve orders with specific field schedules.
- **Role-Based Access Control**: Secure endpoints using JWT authentication and role-based permissions (Admin & Customer).
- **Graceful Shutdown**: Robust handling of HTTP server and Kafka consumers to ensure zero data loss during deployments.
- **Kafka Integration**: Asynchronous payment status updates via Kafka consumer.
- **Rate Limiting**: Integrated protection against brute-force and spam using `tollbooth`.
- **Database Migrations**: Automatic schema migration using GORM.

## Tech Stack

- **Language**: Go (Golang)
- **Framework**: [Gin Gonic](https://github.com/gin-gonic/gin)
- **Database**: PostgreSQL with [GORM](https://gorm.io/)
- **Message Broker**: Apache Kafka ([Sarama](https://github.com/IBM/sarama))
- **CLI Tool**: [Cobra](https://github.com/spf13/cobra)
- **Configuration**: Viper & Godotenv
- **Logging**: Logrus

## Prerequisites

- Go 1.26 or higher
- PostgreSQL
- Apache Kafka
- Docker & Docker Compose (optional)

##  Getting Started

### 1. Clone the Repository
```bash
git clone <repository-url>
cd order-service
```

### 2. Environment Configuration
Copy the `.env.example` file to `.env` and fill in your credentials:
```bash
cp .env.example .env
```

### 3. Install Dependencies
```bash
go mod download
```

### 4. Running the Application

#### Local Development
```bash
go run main.go serve
# OR using Makefile
make build
./order-service serve
```

#### Using Docker
```bash
make docker-compose
```

## API Endpoints

All endpoints are prefixed with `/api/v1`.

| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| `POST` | `/order` | Customer | Create a new order |
| `GET` | `/order` | Admin, Customer | List all orders (Paginated) |
| `GET` | `/order/:uuid` | Admin, Customer | Get order details by UUID |
| `GET` | `/order/user` | Customer | Get order history for authenticated user |

### Query Parameters for Pagination (`GET /order`)
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10)
- `sort_column`: Column to sort by (`id`, `amount`, `status`, etc.)
- `sort_order`: `asc` or `desc`

## Project Structure

- `cmd/`: Application entry point and CLI commands.
- `controllers/`: HTTP handlers and Kafka message handlers.
- `services/`: Business logic layer.
- `repositories/`: Database abstraction layer.
- `domain/`: Data structures (DTOs and Models).
- `clients/`: External service clients (User, Payment, Field).
- `middlewares/`: Gin middlewares (Auth, Role Check, Rate Limit).
- `constants/`: Global constants and error definitions.
- `common/`: Utility functions and shared response helpers.

## Error Handling

The service uses a centralized error mapping system. Errors from external services (like User Service) are mapped to internal domain errors to provide accurate HTTP status codes:

- `400 Bad Request`: Invalid input or UUID format.
- `401 Unauthorized`: Missing or invalid JWT token.
- `403 Forbidden`: Insufficient permissions for the role.
- `404 Not Found`: Resource (Order/User/Field) does not exist.
- `409 Conflict`: Field is already booked.
- `500 Internal Server Error`: Unexpected server issues.

## Kafka Consumers

The service listens to the following topics:
- `payment-status-updated`: Updates order status based on payment events (Settlement, Expired, Pending).

