# Weather Alert Service

## Overview
Weather Alert Service is a Go-based microservice that provides:
- Retrieving current weather data for cities
- User subscriptions to weather alerts based on custom conditions with email confiramatioin
- Automated daily email notifications when conditions are met

This service uses Gin for HTTP handling, GORM for MySQL interactions, and Google Wire for dependency injection.

## Architecture

```
┌─────────┐      ┌──────────────┐      ┌───────────┐
│  Client │<---->│ HTTP Router  │<---->│ Controllers│
└─────────┘      └──────────────┘      └───────────┘
                                     ↕           ↕
                              ┌───────────┐   ┌──────────┐
                              │ Services  │   │ Repos     │
                              └───────────┘   └──────────┘
                                     ↕           ↕
                              ┌───────────┐   ┌──────────┐
                              │  GORM DB  │   │ MailHog   │
                              └───────────┘   └──────────┘
```

### Layers
- **cmd/app**: Entry point, bootstraps DI, starts scheduler and HTTP server
- **internal/http**:
  - **controllers**: HTTP handlers, response formatting
  - **routes**: Route registration with DI
- **pkg**:
  - **config**: Environment loading
  - **database**: MySQL connection and migrations
  - **models**: GORM models for Weather and Subscription
  - **repository**: Interfaces and GORM-based implementations
  - **services**: Business logic (weather retrieval, subscription management, notification evaluation)
  - **utils**: Email sending utility
- **internal/scheduler**: Cron job for daily alert checks
- **wire.go**: Google Wire setup for dependency injection

## Key Architectural Decisions
- **Dependency Injection (Google Wire)**: Ensures loose coupling, easier testing, and clear wiring of dependencies in `InitializeApp()`.
- **Gin Framework**: Lightweight, efficient HTTP router with built-in middleware support.
- **GORM ORM**: Simplifies database operations and migrations.
- **Repository Pattern**: Abstracts data access, facilitating mocking in unit tests.
- **Layered Structure**: Separates concerns for controllers, services, repositories, and database logic.

## Technologies

- **Language:** Go 1.24+
- **Web Framework:** gin‑gonic/gin
- **HTTP Validator:** go-playground/validator/v10
- **Environment Loader:** joho/godotenv
- **ORM:** GORM (`gorm.io/gorm`)
- **SQL Driver:** GORM MySQL Driver (`gorm.io/driver/mysql`)
- **Database:** MySQL (production), SQLite (integration tests)
- **Dependency Injection:** Google Wire (`github.com/google/wire`)
- **Scheduler:** robfig/cron (`github.com/robfig/cron/v3`)
- **Email:** gomail (`gopkg.in/gomail.v2`) + MailHog (or SMTP)
- **Testing:** testify (`github.com/stretchr/testify`)
- **Containerization:** Docker & Docker Compose
- **API Documentation:** swaggo/swag (`github.com/swaggo/swag`)

## Running Locally

### Prerequisites
- Go 1.24+
- Docker & Docker Compose
- MySQL (or use Docker Compose)

### Environment Variables
Create a `.env` file in project root:
```ini
DB_USER=root
DB_PASS=
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=weatheralertservicebd
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=youremail@example.com
SMTP_PASS=yourpassword
```

### Docker Compose
```bash
docker-compose up -d
docker-compose run app wire  # generate wire_gen.go
go run cmd/app/main.go
```

### Manual MySQL Setup
```sql
CREATE DATABASE weatheralertservicebd;
```

### Run Application
```bash
go install github.com/google/wire/cmd/wire@latest
wire                # generate DI code
go run cmd/app/main.go
```

## API Endpoints

| Method | Endpoint                         | Description                                     |
|--------|----------------------------------|-------------------------------------------------|
| GET    | `/weather?city={city}`           | Get current weather for a city                  |
| POST   | `/weather`                       | Create or update weather data (`Weather` JSON)  |
| PUT    | `/weather/{city}`                | Update existing weather by city                 |
| POST   | `/subscriptions`                 | Create a subscription                           |
| GET    | `/subscriptions/confirm?token=`  | Confirm email subscription                      |

### Example JSON
**POST /weather**
```json
{
  "city": "Kyiv",
  "temperature": 12.5,
  "humidity": 60,
  "condition": "Clear"
}
```

**POST /subscriptions**
```json
{
  "email": "user@example.com",
  "city": "Kyiv",
  "condition": "temp<0"
}
```

## Testing
- **Unit Tests**: `go test ./pkg/services` covers business logic
- **Integration Tests**: `go test ./tests/integration` exercises full HTTP API

```bash
go test ./pkg/services -v
go test ./tests/integration -v
```

## Swagger / OpenAPI (Optional)
You can integrate Swagger UI by generating OpenAPI spec with `swaggo/swag`, then serve at `/docs`.

---
Feel free to explore and extend the service according to your needs!



