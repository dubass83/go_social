# GoSocial

A social media API and web application built as a learning project. The backend is a RESTful JSON API written in Go; the frontend is a React single-page application.

## Stack

| Layer | Technology |
|-------|-----------|
| API | Go 1.24, Chi router, JWT auth |
| Database | PostgreSQL 16 with golang-migrate |
| Cache | Redis 6 (optional) |
| Frontend | React 19, TypeScript, Vite, React Router v7 |
| Dev reload | Air (backend), Vite HMR (frontend) |

---

## Prerequisites

- Go 1.24+
- Node.js 20+
- Docker & Docker Compose (for PostgreSQL and Redis)
- [Air](https://github.com/air-verse/air) – hot reload: `go install github.com/air-verse/air@latest`
- [golang-migrate CLI](https://github.com/golang-migrate/migrate) – migrations: `brew install golang-migrate`

---

## Quick Start

### 1. Configure environment

Copy the example and fill in values:

```bash
cp .env.example .env   # edit DB_ADDR, JWT_SECRET, MAIL_* etc.
```

Minimum `.env` for local development:

```env
DB_ADDR=postgres://postgres:password@localhost:5432/social?sslmode=disable
JWT_SECRET=changeme
BASIC_AUTH_USERNAME=admin
BASIC_AUTH_PASSWORD=secret
```

### 2. Start the API

```bash
make run        # starts DB, runs migrations, then starts the API
```

The API is available at `http://localhost:8080`.

### 3. Start the web UI

```bash
make web        # starts the Vite dev server
```

The UI is available at `http://localhost:5173`.

---

## Make Targets

| Target | Description |
|--------|-------------|
| `make run` | Start DB, run migrations, run API with `go run` |
| `make repl` | Start DB, run migrations, run API with Air (hot reload) |
| `make web` | Start the React frontend dev server |
| `make build` | Compile the API binary to `bin/myapp` |
| `make start` | Build and start the compiled binary |
| `make stop` | Stop the binary and database |
| `make restart` | Stop then start |
| `make test` | Run all tests |
| `make race_test` | Run tests with race detector |
| `make start_db` | Start the PostgreSQL Docker container |
| `make stop_db` | Stop and remove the PostgreSQL container |
| `make start_redis` | Start the Redis Docker container |
| `make stop_redis` | Stop and remove the Redis container |
| `make migrate_up` | Apply all pending migrations |
| `make migrate_down` | Roll back all migrations |
| `make migrate_up1` | Apply the next single migration |
| `make migrate_down1` | Roll back the last migration |
| `make new_migration name=<name>` | Scaffold a new migration file |
| `make seed` | Seed the database with sample data |
| `make gen-docs` | Regenerate Swagger documentation |

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `API_ADDR` | `:8080` | Address the API listens on |
| `EXTERNAL_URL` | `localhost:8080` | Public base URL (used in emails) |
| `FRONTEND_URL` | `http://localhost:5173` | Allowed CORS origin |
| `DB_ADDR` | — | PostgreSQL connection string |
| `DB_MAX_OPEN_CONNS` | `30` | Max open DB connections |
| `DB_MAX_IDLE_CONNS` | `30` | Max idle DB connections |
| `DB_MAX_IDLE_TIME` | `10m` | Max connection idle time |
| `JWT_SECRET` | `veryStrongSecret` | HMAC secret for signing tokens |
| `JWT_EXPIRY` | `3600` | Token lifetime in seconds |
| `BASIC_AUTH_USERNAME` | — | Username for the health endpoint |
| `BASIC_AUTH_PASSWORD` | — | Password for the health endpoint |
| `CACHE_ADDR` | `localhost:6379` | Redis address |
| `CACHE_ENABLE` | `false` | Enable Redis user cache |
| `RATE_LIMIT_REQUESTS` | `100` | Max requests per window |
| `RATE_LIMIT_TIMEFRAME` | `60` | Rate-limit window in seconds |
| `RATE_LIMIT_ENABLE` | `true` | Toggle rate limiting |
| `MAIL_SERVICE` | `mailtrap` | Mail provider |
| `MAIL_SENDER_NAME` | `GO Social` | From name |
| `MAIL_SENDER_EMAIL` | `noreply@go-social.com` | From address |
| `MAIL_LOGIN` | — | Mail provider login |
| `MAIL_PASSWORD` | — | Mail provider password |

---

## API Reference

All endpoints are prefixed with `/v1`. Responses follow the envelope pattern:

```json
{ "data": <payload> }
{ "error": "message" }
```

### Authentication

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/authentication/user` | — | Register a new user |
| `POST` | `/authentication/token` | — | Sign in, receive JWT |

### Users

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/users/me` | Bearer | Get the currently signed-in user |
| `GET` | `/users/{userID}` | Bearer | Get a user by ID |
| `PUT` | `/users/activate/{token}` | — | Activate account via email token |
| `PUT` | `/users/{userID}/follow` | Bearer | Follow a user |
| `PUT` | `/users/{userID}/unfollow` | Bearer | Unfollow a user |
| `GET` | `/users/{userID}/posts` | Bearer | List posts by a user (paginated) |
| `GET` | `/users/feed` | Bearer | Personalized feed (followed users) |

### Posts

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/posts` | — | List all posts (public, paginated) |
| `POST` | `/posts` | Bearer | Create a new post |
| `GET` | `/posts/{postID}` | Bearer | Get a post with its comments |
| `PATCH` | `/posts/{postID}` | Bearer | Update a post (owner or moderator) |
| `DELETE` | `/posts/{postID}` | Bearer | Delete a post (owner or admin) |
| `POST` | `/posts/{postID}/comments` | Bearer | Add a comment to a post |

### Health & Docs

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/health` | Basic | API health check |
| `GET` | `/swagger/*` | — | Swagger UI |

#### Pagination query parameters

Supported by `/posts`, `/users/feed`, and `/users/{userID}/posts`:

| Parameter | Type | Default | Constraint |
|-----------|------|---------|------------|
| `limit` | int | 10 | 1–100 |
| `offset` | int | 0 | ≥ 0 |
| `sort` | string | `desc` | `asc` or `desc` |
| `search` | string | — | max 100 chars (title & content) |
| `tags` | string | — | comma-separated tag list |

---

## Frontend Routes

| Route | Description |
|-------|-------------|
| `/` | Personal feed (requires auth) |
| `/login` | Sign-in form |
| `/register` | Registration form |
| `/confirm/:token` | Account activation (from email link) |
| `/posts/new` | Create a new post |
| `/posts/:postID` | Post detail with comments |
| `/posts/:postID/edit` | Edit an existing post |
| `/users/:userID` | User profile with their posts and follow button |

---

## Project Structure

```
.
├── cmd/
│   ├── api/                  # HTTP handlers, router, middleware
│   └── db/
│       ├── migration/        # SQL migration files
│       └── seed/             # Database seed script
├── internal/
│   ├── store/                # Repository layer (Posts, Users, Comments, …)
│   ├── db/                   # Database connection
│   └── env/                  # Environment variable helpers
├── web/                      # React frontend
│   └── src/
│       ├── api.ts            # Typed API client
│       ├── types.ts          # Shared TypeScript interfaces
│       ├── context/          # AuthContext (JWT + current user)
│       ├── components/       # Layout, PostCard, ProtectedRoute
│       └── pages/            # One file per route
├── Makefile
└── docker-compose.yml
```

---

## User Roles

| Role | Level | Permissions |
|------|-------|-------------|
| `user` | 1 | Create posts and comments, follow others |
| `moderator` | 2 | Edit any post |
| `admin` | 3 | Delete any post |

---

## Registration Flow

1. `POST /v1/authentication/user` — creates the account and sends an activation email.
2. User clicks the link in the email → `PUT /v1/users/activate/{token}` (or the `/confirm/:token` page in the UI).
3. Account is activated; user can now sign in with `POST /v1/authentication/token`.
4. All subsequent requests include the returned JWT as `Authorization: Bearer <token>`.
