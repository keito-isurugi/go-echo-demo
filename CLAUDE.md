# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Essential Development Commands

### Environment Setup
```bash
# Copy environment variables template
cp env.example .env
# Configure .env with database settings and JWT_SECRET_KEY
```

### Database Setup
```bash
# Initialize database with tables and test data
psql -U postgres -d go_echo_demo -f init_db.sql

# Apply RBAC migrations
psql -U postgres -d go_echo_demo -f migrate_rbac.sql
```

### Running the Application
```bash
# Development with hot reload
air

# Direct run
go run cmd/main.go

# Build
go build -o app ./cmd/main.go

# Docker Compose (includes PostgreSQL)
docker-compose up
```

## High-Level Architecture

This project follows **Clean Architecture** with clear separation of concerns:

### Layer Responsibilities
- **Domain** (`internal/domain/`): Business entities and interface definitions. All other layers depend on this.
- **Use Case** (`internal/usecase/`): Business logic implementation. Orchestrates between domain and repository.
- **Repository** (`internal/repository/`): Data persistence interfaces and implementations.
- **Handler** (`internal/handler/`): HTTP request handling, split into:
  - `api/`: REST API endpoints returning JSON
  - `frontend/`: Server-side rendered HTML pages
- **Middleware** (`internal/middleware/`): Cross-cutting concerns (auth, CORS, logging)
- **Infrastructure** (`internal/infrastructure/`): External service configurations (OAuth, Casbin)

### Dependency Flow
All dependencies are injected in `cmd/main.go` following this pattern:
```
Repository → UseCase → Handler → Router
```

## Authentication System

The application supports multiple authentication methods:

### JWT Authentication
- Login endpoint: `POST /api/auth/login`
- Protected routes require `Authorization: Bearer <token>` header
- Token validation middleware: `internal/middleware/jwt.go`
- Tokens expire after 24 hours

### OAuth Providers
- Configured in `internal/infrastructure/oauth.go`
- Flow: `/auth/{provider}` → External auth → `/auth/{provider}/callback`
- State parameter used for CSRF protection
- Providers: Google, LINE (extensible for others)

### Basic and Digest Auth
- Middleware implementations for HTTP basic/digest authentication
- Used for specific demo endpoints

## RBAC System

Two RBAC implementations coexist:

### Database-based RBAC
- Tables: `roles`, `permissions`, `user_roles`, `role_permissions`
- Permission format: `resource:action` (e.g., "users:read")
- Managed through `/rbac-admin` interface

### Casbin RBAC
- Model: `config/rbac_model.conf`
- Policies: `config/rbac_policy.csv`
- Managed through `/casbin-admin` interface
- More flexible policy-based access control

## Database Schema

Key tables and their relationships:
- `users`: Core user data with OAuth provider support
- `roles`: Predefined roles (admin, user, guest)
- `permissions`: Granular permissions
- Association tables for many-to-many relationships

## OAuth Provider Extension

To add a new OAuth provider:
1. Implement the `OAuthProvider` interface in domain
2. Create provider implementation in infrastructure
3. Register in `NewOAuthProviderFactory()`
4. Add frontend button in login template

## Test Credentials
- Admin: `user1@example.com` / `password123`
- User: `user2@example.com` / `password456`

## Important Patterns

1. **Error Handling**: Errors bubble up from repository → usecase → handler
2. **Transaction Management**: Handled at repository level
3. **Middleware Chaining**: Applied in route definitions
4. **Template Inheritance**: Base templates with header/footer
5. **Static File Serving**: `/static/` directory for CSS and assets