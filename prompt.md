# Prompt AI Agent — Build Management Siswa API

You are a senior backend engineer responsible for building a production-grade REST API for a system called **Management Siswa**.

The project structure already exists and must be respected exactly.

---

# Tech Stack

Use the following stack strictly:

* Go 1.24+
* Fiber v2
* PostgreSQL via NeonDB
* SQLX (NOT ORM)
* Argon2 password hashing
* JWT Authentication
* Role Based Access Control (RBAC)
* Swagger / OpenAPI
* Validator
* Zerolog
* UUID
* golang-migrate
* godotenv

Do NOT use:

* GORM
* MongoDB
* ORM abstraction layer
* GraphQL
* Repository generators
* Heavy frameworks

---

# Architecture Rules

Follow clean architecture with domain separation.

Project structure already exists and MUST be followed.

@project_structure.text

```text
management-siswa-api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   ├── database/
│   │   ├── migrations/
│   │   └── seeds/
│   ├── domain/
│   │   ├── auth/
│   │   ├── student/
│   │   ├── progress/
│   │   └── user/
│   ├── middleware/
│   ├── routes/
│   ├── constants/
│   └── dto/
├── pkg/
│   ├── utils/
│   └── token/
├── docs/
├── .env
├── .env.example
├── go.mod
└── go.sum
```

---

# Database

Database = PostgreSQL (NeonDB)

Use UUID as primary key.

Use SQLX for queries.

Create migrations for all tables.

---

# ERD

@ERD.mermaid

---

# API Requirements

Implement all API endpoints below @apispect.json

# Required Middleware

Implement middleware:

### JWT Middleware

* validate JWT
* inject user context

### RBAC Middleware

Support:

* admin
* mentor

### Request ID Middleware

### Logger Middleware

### Recover Middleware

### Rate Limiter

Apply limiter to:

```text
/auth/login
```

---

# Validation

Use validator package.

Return consistent validation errors.

Example:

```json
{
  "success": false,
  "message": "validation error",
  "errors": {
    "email": "invalid email"
  }
}
```

---

# Response Format

Use consistent JSON response.

Success:

```json
{
  "success": true,
  "message": "success",
  "data": {}
}
```

Error:

```json
{
  "success": false,
  "message": "error message"
}
```

---

# Logging

Use Zerolog.

Log:

* request method
* path
* latency
* error
* request ID

---

# Database Migration

Create SQL migration files for:

* users
* students
* progress

Use:

```text
internal/database/migrations
```

---

# Security Rules

* Never expose password_hash
* Always hash passwords with Argon2
* JWT secret loaded from .env
* Use environment config
* Use prepared queries
* Protect private routes

---

# Swagger

Generate Swagger documentation.

Swagger endpoint:

```text
/docs
```

---

# Code Quality Rules

* Keep handlers thin
* Business logic only in service layer
* Repository only for database queries
* DTO separate from models
* No duplicated logic
* No fat handlers
* No SQL inside handlers
* Use interfaces for services and repositories
* Use dependency injection

---

# Deliverables

Generate:

1. Working REST API
2. Migration SQL
3. DTOs
4. Models
5. Middleware
6. Auth system
7. Swagger
8. Route registration
9. Database connection
10. Env loader
11. Logger setup
12. RBAC
13. JWT helper
14. Pagination helper
15. Response helper
16. Validation helper
17. README installation guide
18. Example .env

---

# Important

Do NOT generate pseudo-code.

Generate production-grade code.

Code must compile.

Code must follow clean architecture.

Do not skip any layer.

Generate complete backend implementation.

# Execution Order

1. Read project structure
2. Read ERD
3. Read API specification
4. Generate migrations
5. Generate models
6. Generate repositories
7. Generate services
8. Generate handlers
9. Generate middleware
10. Generate routes
11. Generate swagger
12. Generate README
