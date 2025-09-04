# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based social media API application using Chi router, PostgreSQL database, and database migrations. The application implements a clean architecture pattern with separation between API handlers, business logic, and data storage layers.

## Development Commands

All commands use the provided Makefile:

### Build and Run
- `make build` - Build the binary to `bin/myapp`
- `make run` - Start database, run migrations, and run the application with hot reload
- `make start` - Build and start the compiled application
- `make stop` - Stop the running application and database
- `make restart` - Stop and start the application
- `make repl` - Start with Air for hot reload development

### Testing
- `make test` - Run all tests (includes database startup and migrations)
- `make race_test` - Run tests with race condition detection

### Database Management
- `make start_db` - Start PostgreSQL database container and create database
- `make stop_db` - Stop and remove database container
- `make migrate_up` - Run all pending migrations
- `make migrate_up1` - Run next single migration
- `make migrate_down` - Rollback all migrations
- `make migrate_down1` - Rollback last migration
- `make new_migration name=<migration_name>` - Create new migration files

### Redis (if needed)
- `make start_redis` - Start Redis container
- `make stop_redis` - Stop Redis container

## Architecture

### Directory Structure
- `cmd/api/` - Application entry point and HTTP handlers
- `cmd/db/migration/` - Database migration files
- `internal/store/` - Data access layer with repository pattern
- `internal/db/` - Database connection logic
- `internal/env/` - Environment variable utilities

### Key Components
- **API Layer** (`cmd/api/`): Chi-based HTTP router with RESTful endpoints
- **Storage Layer** (`internal/store/`): Repository pattern with interfaces for Posts, Users, and Comments
- **Database Layer** (`internal/db/`): PostgreSQL connection with connection pooling
- **Migration System**: Uses golang-migrate for database schema management

### API Structure
- Base endpoint: `/v1/`
- Health check: `GET /v1/health`
- Posts CRUD: `/v1/posts/` with full CRUD operations
- Middleware: Request ID, logging, recovery, timeout (60s)
- Optimistic locking implemented for post updates using version field

### Database Schema
- Posts table with optimistic locking via version column
- Users table for authentication/user management
- Comments table linked to posts
- Uses PostgreSQL array types for tags

### Key Features
- JSON response envelope structure for consistent API responses
- Context middleware for post operations
- Database connection pooling with configurable limits
- Structured logging with zerolog
- Request validation using go-playground/validator

## Environment Variables
- `API_ADDR` - Server address (default: `:8080`)
- `DB_ADDR` - Database connection string
- `DB_MAX_OPEN_CONNS` - Max open database connections (default: 30)
- `DB_MAX_IDLE_CONNS` - Max idle database connections (default: 30)
- `DB_MAX_IDLE_TIME` - Max idle time for connections (default: `10m`)

## Docker Services
- PostgreSQL 16.3 on port 5432
- Redis 6.2 on port 6379 (available but not actively used)

## Development Notes
- Application uses Go 1.24.6
- Hot reload development with Air (via `make repl`)
- Database operations require running `make start_db` first
- All tests require database to be running with migrations applied
- Use Angular commit style for commit messages (e.g., `feat:`, `fix:`, `docs:`, `refactor:`)
