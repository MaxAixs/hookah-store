# AGENTS.md

# Project Overview

This project is a scalable e-commerce platform built using microservices
architecture.

Repository name:

hookah-store

All microservices are developed inside a single monorepository.

# Technology Stack

-   Golang 1.22+
-   Gin HTTP framework
-   gRPC for internal communication
-   PostgreSQL
-   Apache Kafka
-   Docker
-   Docker Compose
-   Grafana
-   Prometheus

# Repository Structure

    hookah-store/

    ├── user-service/
    ├── product-service/
    ├── cart-service/
    ├── order-service/
    ├── payment-service/
    ├── notification-service/

    ├── docker-compose.yml
    ├── Makefile
    └── README.md

Each microservice must contain:

    service-name/
    ├── cmd/
    ├── internal/
    ├── pkg/
    ├── migrations/
    ├── configs/
    ├── Dockerfile
    ├── go.mod
    └── README.md

# Architecture Rules

Use Clean Architecture:

Handler -\> Service -\> Repository -\> Database

## Handler Layer

Responsible for: - HTTP request parsing - Validation - Calling use
cases - HTTP responses

Handlers must not contain business logic.

## Use Case Layer

Responsible for: - Business logic - Domain rules - Application workflows

Must not depend on HTTP or database implementation.

## Repository Layer

Responsible only for: - PostgreSQL queries - Data persistence

Use interfaces between layers.

### Repository separation

Each repository struct must be split by domain concern, each in its own file:

-   `auth.go` — `AuthRepo` implements `AuthRepository` (auth operations: `GetByEmail`, `UpdatePassword`)
-   `user.go` — `UserRepo` implements `UserRepository` (user CRUD: `Create`, `GetByID`, `Update`, `Delete`)

This is a general architectural rule: every repository has its own zone of responsibility. Do not mix concerns in one file.

Repository methods must be on the correct receiver type (`*AuthRepo` or `*UserRepo`), not mixed.

Interfaces are defined in `internal/repository/interface.go`.

Services depend on interfaces, not concrete types. Each service receives only the repositories it needs:

-   `AuthService` — needs both `AuthRepository` (sign-in, password reset) and `UserRepository` (sign-up creates user)
-   `AdminService` — needs only `UserRepository`

# Dependency Injection

Always use constructors.

Avoid: - global variables - hidden dependencies - singleton state

# Context Rules

Every function below handler level must accept context.Context.

Never create context.Background() inside business logic.

# Logging Rules

Use Go standard structured logging:

-   log/slog

Do not use third-party logging libraries.

Every function below handler level must have a function context
constant.

Format:

{service}*{layer}*{action}

Example:

``` go
const fc = "user-service_service_create"
```

Logs must include: - fc - service name - request ID - error details

Use structured JSON logs.

# Error Handling

Always wrap errors:

``` go
return fmt.Errorf("create user: %w", err)
```

Never ignore errors.

# HTTP API

Use Gin framework.

Handlers must only: - parse request - validate input - call service -
return response

# gRPC

Internal service communication uses gRPC.

All protobuf files are stored in:

    proto/api/

Generated code must not be manually edited.

# Kafka

Kafka is used for asynchronous communication.

Events examples:

-   user.created
-   order.created
-   payment.completed
-   notification.sent

Events should contain: - event_id - event name - timestamp - producer -
payload

# PostgreSQL

Each microservice owns its database.

Do not share tables between services.

All schema changes require migrations.

Never modify existing migrations.

# Docker

Every service must have a Dockerfile.

Local development uses Docker Compose.

No secrets inside Dockerfiles.

Use environment variables.

# Testing

Use: - table-driven tests - mocks - integration tests

Run:

go test ./...

# Code Quality

Required:

go fmt ./... go vet ./... golangci-lint run

# Observability

Use:

-   Prometheus
-   Grafana

Track: - request latency - request count - errors - database errors -
Kafka errors

# Security

Rules: - validate external input - never log secrets - secure JWT
configuration - use timeouts - use retries - use rate limiting when
required

# Local Development

The project is developed as a monorepository.

Communication:

-   gRPC for synchronous operations
-   Kafka for asynchronous events

Start locally:

docker compose up

# Additional Components

In addition to the core microservices, the platform will include the following
components to enhance scalability, reliability, and manageability:

## API Gateway

Serves as the entry point for all client requests, routing them to the
appropriate microservice. Handles cross-cutting concerns:

-   JWT authentication and token validation
-   Rate limiting
-   Request routing and load balancing
-   Common CORS, logging, tracing headers

Consider: Kong, Traefik, or NGINX.

All JWT validation happens at the gateway level. Downstream services receive
`X-User-ID` / `X-User-Email` headers and do not need to import JWT libraries.

## Service Discovery

Automatically detects and manages service instances.

Consider: Consul or Eureka.

For local development Docker Compose DNS (service name resolution) is
sufficient.

## Centralized Logging

Aggregates logs from all microservices for easy monitoring and debugging.

Consider: ELK stack (Elasticsearch, Logstash, Kibana).

## Docker & Docker Compose

Containerize each microservice and manage their orchestration, networking, and
scaling.

Docker Compose is used to define and manage multi-container applications
locally.

## CI/CD Pipeline

Automates the build, test, and deployment process of each microservice.

Consider: Jenkins, GitLab CI, or GitHub Actions.

# AI Agent Rules

When modifying code:

-   follow existing architecture
-   prefer interfaces
-   keep changes minimal
-   explain architectural decisions
-   do not introduce new technologies without discussion
-   ASK before making any changes — do not modify files or create new ones without explicit user approval

Goal:

Build a production-ready scalable Go microservices e-commerce platform.
