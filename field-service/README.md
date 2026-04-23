# Field Service

A microservice for managing booking system fields, their operational times, and booking schedules. It handles field data, schedule generation, image uploads via Cloudflare R2, and interacts with a user service for role-based access control.

##  Tech Stack

- **Language**: Go
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL (with GORM)
- **Storage**: Cloudflare R2
- **Containerization**: Docker & Docker Compose
- **Live Reload**: Air

##  Prerequisites

- Go 1.26+
- PostgreSQL
- Docker & Docker Compose
- Make

##  Configuration

Copy the `.env.example` to `.env` and fill in your credentials:

```bash
cp .env.example .env
```

### Key Environment Variables

- `PORT`: Application port (default: 8002)
- `SIGNATURE_KEY`: Secret key for API validation
- `DB_*`: PostgreSQL connection details
- `CLOUDFLARE_*` / `R2_*`: Cloudflare R2 credentials for image storage
- `RATE_LIMITER_*`: Settings for request rate limiting

## API Endpoints

### Fields (`/api/v1/field`)

| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| GET | `/` | Public | Get all fields without pagination |
| GET | `/:uuid` | Public | Get field details by UUID |
| GET | `/pagination` | Admin/Customer | Get fields with pagination |
| POST | `/` | Admin | Create a new field |
| PUT | `/:uuid` | Admin | Update field details |
| DELETE | `/:uuid` | Admin | Delete a field |

### Field Schedules (`/api/v1/field/schedule`)

| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| GET | `/lists/:uuid` | Public | Get schedules by Field ID and Date |
| PATCH | `/` | Public/Service | Update schedule status |
| GET | `/pagination` | Admin/Customer | Get schedules with pagination |
| GET | `/:uuid` | Admin/Customer | Get schedule details by UUID |
| POST | `/` | Admin | Create a specific schedule |
| POST | `/one-month` | Admin | Generate daily schedules for a month |
| PUT | `/:uuid` | Admin | Update schedule details |
| DELETE | `/:uuid` | Admin | Delete a schedule |

### Times (`/api/v1/time`)

| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| GET | `/` | Admin | Get all operational times |
| GET | `/:uuid` | Admin | Get time details by UUID |
| POST | `/` | Admin | Create a new operational time |

## Running the Application

### Development (Local)

1.  **Prepare tools**:
    ```bash
    make watch-prepare
    ```
2.  **Run with hot reload**:
    ```bash
    make watch
    ```

### Docker

**Run service and database**:
```bash
make docker-compose
```

Alternatively, build the image manually:
```bash
make docker-build tag=1.0.0
```
