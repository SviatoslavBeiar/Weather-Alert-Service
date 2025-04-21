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
  - **services**: Business logic (weather retrieval, subscription management, notification evaluation, unit tests)
  - **utils**: Email sending utility
- **internal/scheduler**: Cron job for daily alert checks
- **wire.go**: Google Wire setup for dependency injection

## Key Architectural Decisions
- **Dependency Injection (Google Wire)**: Ensures loose coupling, easier testing, and clear wiring of dependencies in `InitializeApp()`.
- **Gin Framework**: Lightweight, efficient HTTP router with built-in middleware support.
- **GORM ORM**: Simplifies database operations and migrations.
- **Repository Pattern**: Abstracts data access, facilitating mocking in unit tests.
- **Layered Structure**: Separates concerns for controllers, services, repositories, and database logic.

## 🚀 Technologies

- **Language & Version**
  - Go 1.24.2

- **Web Framework & Routing**
  - [Gin](https://github.com/gin-gonic/gin) (v1.10.0)

- **Validation**
  - [go-playground/validator](https://github.com/go-playground/validator) (v10.26.0)

- **Dependency Injection**
  - [Google Wire](https://github.com/google/wire) (v0.6.0)

- **Configuration**
  - [godotenv](https://github.com/joho/godotenv) (v1.5.1)

- **Scheduling**
  - [robfig/cron/v3](https://github.com/robfig/cron) (v3.0.1)

- **Logging**
  - [uber-go/zap](https://github.com/uber-go/zap) (v1.27.0)

- **Email**
  - [gomail](https://github.com/go-gomail/gomail) (v2.0.0)

- **ORM & Database Drivers**
  - [GORM](https://gorm.io) (v1.25.12)
    - MySQL driver (v1.5.7)
    - SQLite driver (v1.5.7)

- **Testing & Mocks**
  - [stretchr/testify](https://github.com/stretchr/testify) (v1.9.0)

- **Indirect Dependencies** (selected)
  - `github.com/go-sql-driver/mysql` (MySQL driver for GORM)
  - `github.com/mattn/go-sqlite3` (SQLite driver)
  - various others pulled in by GORM, Gin, Wire, etc.

> See **go.mod** for the full list of modules and versions.


## Running Locally

### Prerequisites
- Go 1.24+
- Docker & Docker Compose
- MySQL (or use Docker Compose)

### Environment Variables
Create a `.env` file in project root(already created):
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
git clone https://github.com/SviatoslavBeiar/Weather-Alert-Service.git
cd Weather-Alert-Service
docker-compose up -d
```
#### MySQL available at localhost:3306

#### MailHog UI available at http://localhost:8025

#### App  Example http://localhost:8080/weather?city=Kyiv if exit
### Manual MySQL 
```sql
CREATE DATABASE weatheralertservicebd;
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
## Testing Scenarios

| #  | Scenario                                                     | Precondition / Setup                                                                                                                                                    | Trigger / Input                                                                                           | Expected Outcome                                                                                             | Example Email Payload                                                                                                                      |
|----|--------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------|
| 1  | Verified subscription + condition matches → send alert        | **DB:**<br/>  • Weather: `{"city":"Lviv","temperature":4,"humidity":80,"condition":"Clear"}`<br/>  • Subscription: `{"email":"alice@example.com","city":"Lviv","condition":"temp<5","verified":true}` | Cron job runs → calls `EvaluateAndNotify(sub, weather)`                                                    | • Email sent to `alice@example.com`<br/>• `LastSent` updated                                                   | **To:** alice@example.com<br/>**Subject:** Weather Alert for Lviv<br/>**Body:** Condition temp<5 met: current temp 4.0°C                     |
| 2  | Verified subscription + condition does **not** match → no alert | **DB:**<br/>  • Weather: `{"city":"Lviv","temperature":6,"humidity":80,"condition":"Clear"}`<br/>  • Subscription: same as above                                         | Cron job runs → calls `EvaluateAndNotify(sub, weather)`                                                    | • No email sent<br/>• `LastSent` remains unchanged                                                              | *n/a*                                                                                                                                         |
| 3  | Unverified subscription → never send alert                   | **DB:**<br/>  • Weather: any<br/>  • Subscription: `{"email":"bob@example.com","city":"Kyiv","condition":"rain","verified":false}`                                          | Cron job runs                                                                                              | • No email sent<br/>• `LastSent` remains `nil`                                                                  | *n/a*                                                                                                                                         |


