# Payment Service

A robust, scalable, and professional payment gateway integration service built with Go. This service acts as an intermediary between your core application and payment providers (like Midtrans), handling the entire payment lifecycle from transaction creation to automated invoice generation and asynchronous event notification.

## Overview

This service provides a centralized API to manage financial transactions. It is designed to be industry-agnostic, making it suitable for E-commerce, Booking Systems, SaaS platforms, or any application requiring secure payment processing.

### Key Features
- **Payment Gateway Integration**: Seamless integration with Midtrans Snap for secure payment links.
- **Automated Invoicing**: Dynamically generates professional PDF invoices using HTML templates.
- **Asynchronous Processing**: Utilizes Apache Kafka to publish payment events (Pending, Settlement, Expire) to other microservices.
- **Cloud Storage**: Automated upload of generated invoices to Cloudflare R2 (S3-compatible) storage.
- **Rate Limiting**: Built-in protection against API abuse.
- **Robust Security**: Header-based signature validation and role-based access control.

## Tech Stack

- **Language**: [Go (Golang)](https://golang.org/)
- **Web Framework**: [Gin Gonic](https://gin-gonic.com/)
- **Database**: [PostgreSQL](https://www.postgresql.org/) with [GORM](https://gorm.io/)
- **Message Broker**: [Apache Kafka](https://kafka.apache.org/)
- **PDF Engine**: [wkhtmltopdf](https://wkhtmltopdf.org/)
- **Storage**: [Cloudflare R2](https://www.cloudflare.com/products/r2/)
- **Configuration**: [Viper](https://github.com/spf13/viper) & [Cobra](https://github.com/spf13/cobra)

## Prerequisites

Before running this service, ensure you have the following installed:
- Go 1.26+
- PostgreSQL
- Apache Kafka
- **wkhtmltopdf**: Required for PDF generation.
  ```bash
  # Ubuntu/WSL
  sudo apt install wkhtmltopdf
  ```

##  Environment Variables

Create a `.env` file in the root directory based on `.env.example`:

```env
PORT=8003
APP_ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_NAME=payment_db
DB_USERNAME=postgres
DB_PASSWORD=yourpassword

KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=payment-events

MIDTRANS_SERVER_KEY=your_server_key
MIDTRANS_PRODUCTION=false

R2_ACCESS_KEY_ID=your_key
R2_SECRET_ACCESS_KEY=your_secret
R2_BUCKET_NAME=invoices
R2_PUBLIC_URL=your_public_r2_url
```

##  Getting Started

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd payment-service
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run tests**:
   Ensure all logic is working as expected:
   ```bash
   make test
   ```

4. **Run the service**:
   Using `make` (if available):
   ```bash
   make watch
   ```
   Or directly:
   ```bash
   go run main.go serve
   ```

##  Testing

The project includes comprehensive unit tests for Repositories, Services, and Controllers using `testify` and `sqlmock`.

To run the tests:
```bash
go test ./...
```
Or via Makefile:
```bash
make test
```

##  API Documentation

This project uses [Swaggo](https://github.com/swaggo/swag) to generate API documentation.

To generate/update documentation:
```bash
make swagger
```

Once the service is running, you can access the Swagger UI at:
`http://localhost:8003/swagger/index.html`

##  API Endpoints

### Public Routes (Webhook)
- `POST /api/v1/payments/webhook` - Handles payment notifications from providers.

### Protected Routes (Requires Auth)
- `POST /api/v1/payments` - Create a new payment transaction.
- `GET /api/v1/payments` - List all transactions (Pagination supported).
- `GET /api/v1/payments/:uuid` - Get detailed information about a specific transaction.

## Invoicing Logic

When a payment status changes to `settlement`:
1. The service fetches the transaction details.
2. An HTML template (`template/invoice.html`) is rendered with the data.
3. `wkhtmltopdf` converts the HTML to a PDF.
4. The PDF is uploaded to Cloudflare R2.
5. The `invoice_link` is updated in the database and published to Kafka.
