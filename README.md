# Weather Alert Service
## Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [Architecture](#architecture)
- [Layers](#layers)
- [Project Directory Structure](#project-directory-structure)
- [Key Architectural Decisions](#key-architectural-decisions)
- [🚀 Technologies](#-technologies)
- [Running Locally](#running-locally)
- [Environment Variables](#environment-variables)
- [API Endpoints](#api-endpoints)
- [Testing Scenarios](#testing-scenarios)
# Overview
Weather Alert Service is a Go-based microservice focusing on verified email subscriptions before sending any alerts. Users must confirm their subscription via email; only then will they receive automated notifications when defined weather conditions are met.

## Key Features
### Current Weather Retrieval
- Fetch the latest weather data (temperature, humidity, sky condition) for any city.

### Email‑Confirmed Subscriptions
 - Users subscribe with a custom condition (e.g., temp<0), receive a confirmation email, and only verified email addresses will be alerted.

### Automated Alerts
- A daily cron job evaluates registered conditions and sends alerts only for verified subscriptions.

### This service uses Gin for HTTP handling, GORM for MySQL interactions, and Google Wire for dependency injection.

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
### Project Directory Structure

```
Weather-Alert-Service/
├── cmd/
│   └── app/
│       └── main.go         # Entry point, bootstraps DI, starts scheduler and HTTP server
├── internal/
│   ├── http/
│   │   ├── controllers/    # HTTP handlers (controllers)
│   │   └── routes/         # Route registration with DI
│   └── scheduler/          # Cron job for daily alert checks
├── pkg/
│   ├── config/             # Environment loading (Config struct)
│   ├── database/           # MySQL connection and migrations
│   ├── models/             # GORM models for Weather and Subscription
│   ├── repository/         # Interfaces and GORM-based implementations
│   ├── services/           # Business logic (weather retrieval, subscription management, notifications, unit tests)
│   ├── utils/              # Email sending utility, error helpers
│   └── validation/         # Custom validators for request binding
├── app/                    # Google Wire setup and InitializeApp
│   └── wire.go             # DI definitions
├── wire.go                 # (Alternative root DI definitions, may be removed)
├── docker-compose.yml      # Docker Compose for MySQL and MailHog
├── .env.example            # Sample environment variables
├── go.mod
├── go.sum
└── README.md               # Project overview, setup, usage
```



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
SMTP_USER=
SMTP_PASS=
```

### Docker Compose / Aap, DB, MailHog
#### DO not fogert run docker
```bash
git clone https://github.com/SviatoslavBeiar/Weather-Alert-Service.git
cd Weather-Alert-Service
docker-compose up -d
```
#### MySQL available at localhost:3306

#### MailHog UI available at http://localhost:8025

#### App weatheralertservicebd Example req http://localhost:8080/weather?city=Kyiv if exit
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
  "temperature": 1,
  "humidity": 60,
  "condition": "Clear"
}
```

**POST /subscriptions**
```json
{
  "email": "user@example.com",
  "city": "Kyiv",
  "condition": "temp<2"
}
```
## Testing Scenarios
#### Confirm email via MailHog
| #  | Scenario                                                     | Precondition / Setup                                                                                                                                                    | Trigger / Input                                                                                           | Expected Outcome                                                                                             | Example Email Payload                                                                                                                      |
|----|--------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------|
| 1  | Verified subscription + condition matches → send alert        | **DB:**<br/>  • Weather: `{"city":"Lviv","temperature":4,"humidity":80,"condition":"Clear"}`<br/>  • Subscription: `{"email":"alice@example.com","city":"Lviv","condition":"temp<5"}` | Cron job runs → calls `EvaluateAndNotify(sub, weather)`                                                    | • Email sent to `alice@example.com`<br/>• `LastSent` updated                                                   | **To:** alice@example.com<br/>**Subject:** Weather Alert for Lviv<br/>**Body:** Condition temp<5 met: current temp 4.0°C                     |
| 2  | Verified subscription + condition does **not** match → no alert | **DB:**<br/>  • Weather: `{"city":"Lviv","temperature":6,"humidity":80,"condition":"Clear"}`<br/>  • Subscription: same as above                                         | Cron job runs → calls `EvaluateAndNotify(sub, weather)`                                                    | • No email sent<br/>• `LastSent` remains unchanged                                                              | *n/a*                                                                                                                                         |
| 3  | Unverified subscription → never send alert                   | **DB:**<br/>  • Weather: any<br/>  • Subscription: `{"email":"bob@example.com","city":"Kyiv","condition":"Clear"}`                                          | Cron job runs                                                                                              | • No email sent<br/>• `LastSent` remains `nil`                                                                  | *n/a*                                                                                                                                         |


